# Firehose Networks

A wrapper around the [Golang Networks Registry library](https://github.com/pinax-network/graph-networks-libs/tree/main/packages/golang) for usage within Substreams/Firehose products. This library provides additional helpers and utilities commonly used across different projects, as well as support for custom networks not found in the registry for those in development.

## Overview

This library serves as an enhanced interface to The Graph's Networks Registry, specifically tailored for Substreams and Firehose ecosystem needs. It extends the base functionality with:

- **Helper utilities** commonly used across Substreams/Firehose projects
- **Custom network definitions** for development and testing environments
- **Fallback registry support** for offline or restricted environments
- **Enhanced network configuration** with additional metadata

## Features

- **Network Registry Integration**: Seamless integration with The Graph's official Networks Registry
- **Custom Networks**: Support for adding custom network configurations not available in the upstream registry
- **Fallback Mechanism**: Built-in fallback to a local registry when the remote registry is unavailable
- **Helper Functions**: Common utilities for network identification, configuration parsing, and validation
- **Development Support**: Easy addition of test networks and development chains

## Usage

```go
import "github.com/streamingfast/firehose-networks"

// Find a network by ID, alias, or name
network := networks.Find("ethereum-mainnet")
if network != nil {
    fmt.Printf("Network: %s\n", network.FullName)
}

// Get only networks with Substreams endpoints
substreamsNetworks := networks.GetSubstreamsRegistry()
for id, network := range substreamsNetworks {
    fmt.Printf("%s: %v\n", id, network.Services.Substreams)
}

// Get only networks with Firehose endpoints
firehoseNetworks := networks.GetFirehoseRegistry()
for id, network := range firehoseNetworks {
    fmt.Printf("%s: %v\n", id, network.Services.Firehose)
}
```

## Helper Functions

<details>
<summary><strong>Registry Filtering Functions</strong></summary>

### GetSubstreamsRegistry()

Returns a filtered registry containing only networks that have Substreams endpoints configured.

```go
substreamsNetworks := networks.GetSubstreamsRegistry()
for networkID, network := range substreamsNetworks {
    fmt.Printf("Network: %s\n", networkID)
    fmt.Printf("Substreams endpoints: %v\n", network.Services.Substreams)
}
```

This is useful when you need to:
- List all available Substreams-enabled networks
- Validate that a network supports Substreams before attempting connections
- Build network selection UIs for Substreams applications

### GetFirehoseRegistry()

Returns a filtered registry containing only networks that have Firehose endpoints configured.

```go
firehoseNetworks := networks.GetFirehoseRegistry()
for networkID, network := range firehoseNetworks {
    fmt.Printf("Network: %s\n", networkID)
    fmt.Printf("Firehose endpoints: %v\n", network.Services.Firehose)
}
```

This is useful when you need to:
- List all available Firehose-enabled networks
- Validate that a network supports Firehose before attempting connections
- Build network selection UIs for Firehose applications

</details>

<details>
<summary><strong>Network Lookup Functions</strong></summary>

### Find(key string)

Finds a network by ID, alias, full name, or short name. Returns the first match found, with priority given to exact ID matches.

```go
// Find by network ID
network := networks.Find("ethereum-mainnet")

// Find by alias
network = networks.Find("eth")

// Find by full name
network = networks.Find("Ethereum Mainnet")
```

### FindByGenesisBlock(blockNum uint64, blockID string)

Finds a network by matching its genesis block number and hash.

```go
network := networks.FindByGenesisBlock(0, "0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
if network != nil {
    fmt.Printf("Found network: %s\n", network.FullName)
}
```

### FindBySubstreamsEndpoint(endpoint string)

Finds a network that contains the specified Substreams endpoint.

```go
network := networks.FindBySubstreamsEndpoint("https://mainnet.eth.streamingfast.io:443")
if network != nil {
    fmt.Printf("Network: %s\n", network.FullName)
}
```

</details>

<details>
<summary><strong>Configuration Helpers</strong></summary>

### GetBytesEncoding(network *registry.Network)

Returns the bytes encoding format for a given network. Returns `registry.Hex` if no specific encoding is configured.

```go
network := networks.Find("ethereum-mainnet")
encoding := networks.GetBytesEncoding(network)
fmt.Printf("Bytes encoding: %v\n", encoding)
```

This is particularly useful for Firehose applications that need to know how to encode/decode blockchain data for a specific network.

### ScheduleUpdateLatestRegistry(ctx context.Context, interval time.Duration, logger *zap.Logger)

Schedules a background goroutine that periodically updates the registry from the latest remote version at the specified interval. This ensures your application stays up-to-date with the latest network configurations.

```go
import (
    "context"
    "time"
    "go.uber.org/zap"
)

// You are expected to provide ctx and logger
ctx := context.Background()
logger := zap.NewNop()

// Update registry every 30 minutes
networks.ScheduleUpdateLatestRegistry(ctx, 30*time.Minute, logger)
```

Key features:
- **Non-blocking**: Runs in a background goroutine
- **Graceful shutdown**: Respects context cancellation for clean shutdowns
- **Error handling**: Logs errors and continues on failed updates
- **Global update**: Updates the global registry used by all other functions

This is useful for long-running applications that need to stay synchronized with the latest network configurations without manual intervention.

</details>

## Custom Networks

The library supports adding custom network configurations that aren't available in the upstream registry. This is useful for development networks, private chains, or networks not yet in the official registry.

Custom networks are defined in [`overrides.go`](./overrides.go) and automatically merged with the official registry data. See the `TRONMainnet` example in that file for reference on how to structure a custom network definition.

To add your own custom network, follow the same pattern used for the existing overrides.

## Fallback Registry

When the remote registry is unavailable, the library automatically falls back to a local copy stored in `fallback_TheGraphNetworkRegistry_*.json`. This ensures your applications continue to work even in offline environments or when the upstream registry is temporarily unavailable.

## Development

This library is particularly useful for:

- **Substreams developers** who need consistent network configurations
- **Firehose operators** managing multiple blockchain networks
- **Development teams** working with custom or test networks
- **Integration projects** requiring reliable network metadata

## Contributing

When adding new networks or helpers:

1. Ensure compatibility with the upstream Networks Registry format
2. Add appropriate tests for new functionality
3. Update the fallback registry when necessary
4. Document any new helper functions

## License

This project follows the same licensing as the upstream Networks Registry project.