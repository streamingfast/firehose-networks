package networks

import (
	"context"
	_ "embed"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	registry "github.com/pinax-network/graph-networks-libs/packages/golang/lib"
	"go.uber.org/zap"
)

//go:embed fallback_TheGraphNetworkRegistry_0.7.6.json
var embeddedRegistryJSON []byte

// NetworkRegistry is a thin wrapper around a [map[string]*registry.Network] to add some helper methods.
type NetworkRegistry map[string]*registry.Network

var (
	registryNetworks     NetworkRegistry
	registryNetworksOnce sync.Once
)

// getRegistryNetworks fetches and caches all networks from the registry (no filtering).
func getRegistryNetworks() NetworkRegistry {
	registryNetworksOnce.Do(func() {
		reg, err := loadRegistry(registry.FromLatestVersion)
		if err != nil {
			// If the network registry cannot be loaded from the latest version,
			// we launch a Go routine that is going to retry exponentially (with a
			// limit) and update the global registryNetworks variable.
			go backgroundUpdateLatestRegistry(context.Background())

			// Fallback, use embedded JSON
			reg, err = loadRegistry(fromEmbeddedJSON)
			if err != nil {
				panic(fmt.Sprintf("Failed to load registry from both network and embedded JSON: %v", err))
			}
		}

		registryNetworks = reg
	})

	return registryNetworks
}

func loadRegistry(loader func() (*registry.NetworksRegistry, error)) (NetworkRegistry, error) {
	nativeRegistry, err := loader()
	if err != nil {
		return nil, err
	}

	registry := NetworkRegistry{}
	for i, net := range nativeRegistry.Networks {
		registry[net.ID] = &nativeRegistry.Networks[i]
	}

	for _, net := range networkOverrides {
		registry.addCustomNetwork(net, false)
	}

	return registry, nil
}

func fromEmbeddedJSON() (*registry.NetworksRegistry, error) {
	return registry.FromJSON(embeddedRegistryJSON)
}

// GetSubstreamsRegistry returns only networks with Substreams endpoints.
func GetSubstreamsRegistry() NetworkRegistry {
	all := getRegistryNetworks()
	filtered := make(NetworkRegistry)
	for id, net := range all {
		if len(net.Services.Substreams) > 0 {
			filtered[id] = net
		}
	}
	return filtered
}

// GetFirehoseRegistry returns only networks with Firehose endpoints.
func GetFirehoseRegistry() NetworkRegistry {
	all := getRegistryNetworks()
	filtered := make(NetworkRegistry)
	for id, net := range all {
		if len(net.Services.Firehose) > 0 {
			filtered[id] = net
		}
	}
	return filtered
}

// addCustomNetwork can be used to add a custom network to the registry map for testing or development.
func (r NetworkRegistry) addCustomNetwork(network *registry.Network, forced bool) {
	if network == nil || network.ID == "" {
		return // Ignore invalid input
	}

	_, found := r[network.ID]
	if found && !forced {
		// If the network already exists and not forced, we skip adding it.
		return
	}

	r[network.ID] = network
}

// Find returns the network by ID or, if not found, by alias (sorted by network ID), FullName, and ShortName.
func (r NetworkRegistry) Find(key string) *registry.Network {
	if n, ok := r[key]; ok {
		return n
	}
	ids := slices.Collect(maps.Keys(r))
	slices.Sort(ids)
	for _, id := range ids {
		net := r[id]
		if slices.Contains(net.Aliases, key) || net.FullName == key || net.ShortName == key || net.ID == key {
			return net
		}
	}
	return nil
}

// FindByGenesisBlock returns the *registry.Network whose genesis block matches the given blockNum and blockID (hash).
//
// Deprecated: Use FindByFirstStreamableBlock instead, as GenesisBlock has been renamed to FirstStreamableBlock in the network registry.
func (r NetworkRegistry) FindByGenesisBlock(blockNum uint64, blockID string) *registry.Network {
	return r.FindByFirstStreamableBlock(blockNum, blockID)
}

// FindByFirstStreamableBlock returns the *registry.Network whose first streamable block matches the given blockNum and blockID (hash).
func (r NetworkRegistry) FindByFirstStreamableBlock(blockNum uint64, blockID string) *registry.Network {
	for _, network := range r {
		if network.Firehose != nil && network.Firehose.FirstStreamableBlock != nil &&
			uint64(network.Firehose.FirstStreamableBlock.Height) == blockNum &&
			nox(network.Firehose.FirstStreamableBlock.ID) == nox(blockID) {
			return network
		}
	}
	return nil
}

// Find is a shortcut for getRegistryNetworks().Find(key).
func Find(key string) *registry.Network {
	return getRegistryNetworks().Find(key)
}

// GetSubstreamsEndpoint returns the preferred Substreams endpoint for a given network key,
// prioritizing streamingfast.io endpoints when available.
func GetSubstreamsEndpoint(key string) string {
	network := Find(key)
	if network == nil || len(network.Services.Substreams) == 0 {
		return ""
	}

	// First, look for streamingfast.io endpoints
	for _, endpoint := range network.Services.Substreams {
		if strings.Contains(endpoint, "streamingfast.io") {
			return endpoint
		}
	}

	// If no streamingfast.io endpoint found, return the first available endpoint
	return network.Services.Substreams[0]
}

// GetFirehoseEndpoint returns the preferred Firehose endpoint for a given network key,
// prioritizing streamingfast.io endpoints when available.
func GetFirehoseEndpoint(key string) string {
	network := Find(key)
	if network == nil || len(network.Services.Firehose) == 0 {
		return ""
	}

	// First, look for streamingfast.io endpoints
	for _, endpoint := range network.Services.Firehose {
		if strings.Contains(endpoint, "streamingfast.io") {
			return endpoint
		}
	}

	// If no streamingfast.io endpoint found, return the first available endpoint
	return network.Services.Firehose[0]
}

// Returns the bytes encoding for a given network
// Returns the raw BytesEncoding type, Hex if not found.
func GetBytesEncoding(network *registry.Network) registry.BytesEncoding {
	if network != nil && network.Firehose != nil {
		return network.Firehose.BytesEncoding
	}
	return registry.Hex
}

// FindBySubstreamsEndpoint returns the *registry.Network whose Substreams endpoint matches the given endpoint.
func (r NetworkRegistry) FindBySubstreamsEndpoint(endpoint string) *registry.Network {
	for _, net := range r {
		if slices.Contains(net.Services.Substreams, endpoint) {
			return net
		}
	}
	return nil
}

var withInfiniteRetries = backoff.WithMaxTries(0)

func backgroundUpdateLatestRegistry(ctx context.Context) {
	operation := func() (NetworkRegistry, error) {
		return loadRegistry(registry.FromLatestVersion)
	}

	registry, err := backoff.Retry(ctx, operation, withInfiniteRetries, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		// We have been cancelled, nothing to do more
		return
	}

	// We could have used a atomic pointer here, but it's not a big deal,
	// the on the fly update is not expected to be frequent and shouldn't cause
	// any real issues as they are separated instances.
	registryNetworks = registry
}

// ScheduleUpdateLatestRegistry schedules a background update goroutine of the latest registry at the
// specified interval. It runs in a goroutine and updates the global registryNetworks variable. You
// can control it with a context to stop the updates gracefully.
//
// If you don't want any logging, pass nil as the logger parameter.
func ScheduleUpdateLatestRegistry(ctx context.Context, interval time.Duration, logger *zap.Logger) {
	if logger == nil {
		logger = zap.NewNop()
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				// Exit if context is cancelled
				logger.Debug("stopping background registry update due to context cancellation")
				return

			case <-ticker.C:
				registry, err := loadRegistry(registry.FromLatestVersion)
				if err != nil {
					logger.Info("failed to load latest registry, skipping this interval update", zap.Error(err))
					continue
				}

				registryNetworks = registry
			}
		}
	}()
}
