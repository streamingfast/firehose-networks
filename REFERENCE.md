# Firehose Networks - API Reference

## Table of Contents

- [Registry Filtering Functions](#registry-filtering-functions)
  - [GetSubstreamsRegistry()](#getsubstreamsregistry)
  - [GetFirehoseRegistry()](#getfirehoseregistry)
- [Network Lookup Functions](#network-lookup-functions)
  - [Find(key string)](#findkey-string)
  - [FindByFirstStreamableBlock(blockNum uint64, blockID string)](#findbyfirststreamableblockblocknum-uint64-blockid-string)
  - [FindByGenesisBlock(blockNum uint64, blockID string)](#findbygenesisblockblocknum-uint64-blockid-string) *(Deprecated)*
  - [FindBySubstreamsEndpoint(endpoint string)](#findbysubstreamsendpointendpoint-string)
- [Endpoint Helper Functions](#endpoint-helper-functions)
  - [GetSubstreamsEndpoint(key string)](#getsubstreamsendpointkey-string)
  - [GetFirehoseEndpoint(key string)](#getfirehoseendpointkey-string)
- [Configuration Helpers](#configuration-helpers)
  - [GetBytesEncoding(network *registry.Network)](#getbytesencodingnetwork-registrynetwork)
  - [ScheduleUpdateLatestRegistry(ctx context.Context, interval time.Duration, logger *zap.Logger)](#scheduleupdatelatestregistryctx-contextcontext-interval-timeduration-logger-zaplogger)

## Registry Filtering Functions

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

## Network Lookup Functions

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

### FindByFirstStreamableBlock(blockNum uint64, blockID string)

Finds a network by matching its first streamable block number and hash. This is the recommended method for finding networks by block information.

```go
network := networks.FindByFirstStreamableBlock(0, "0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
if network != nil {
    fmt.Printf("Found network: %s\n", network.FullName)
}
```

This function is useful when you have block information from Firehose and need to identify which network it belongs to.

### FindByGenesisBlock(blockNum uint64, blockID string) *(Deprecated)*

**⚠️ Deprecated:** Use `FindByFirstStreamableBlock` instead, as GenesisBlock has been renamed to FirstStreamableBlock in the network registry.

Finds a network by matching its genesis block number and hash. This function is maintained for backward compatibility but internally calls `FindByFirstStreamableBlock`.

```go
// Deprecated - use FindByFirstStreamableBlock instead
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

## Endpoint Helper Functions

### GetSubstreamsEndpoint(key string)

Returns the preferred Substreams endpoint for a given network key, prioritizing streamingfast.io endpoints when available.

```go
endpoint := networks.GetSubstreamsEndpoint("ethereum-mainnet")
if endpoint != "" {
    fmt.Printf("Substreams endpoint: %s\n", endpoint)
    // Use endpoint to connect to Substreams
}
```

This function is useful when you need to:
- Get a ready-to-use Substreams endpoint for a specific network
- Automatically prefer StreamingFast endpoints when multiple options are available
- Handle cases where a network might not have Substreams support (returns empty string)

### GetFirehoseEndpoint(key string)

Returns the preferred Firehose endpoint for a given network key, prioritizing streamingfast.io endpoints when available.

```go
endpoint := networks.GetFirehoseEndpoint("ethereum-mainnet")
if endpoint != "" {
    fmt.Printf("Firehose endpoint: %s\n", endpoint)
    // Use endpoint to connect to Firehose
}
```

This function is useful when you need to:
- Get a ready-to-use Firehose endpoint for a specific network
- Automatically prefer StreamingFast endpoints when multiple options are available
- Handle cases where a network might not have Firehose support (returns empty string)

## Configuration Helpers

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
