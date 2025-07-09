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
import (
    networks "github.com/streamingfast/firehose-networks"
)

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

## API Reference

For detailed documentation of all helper functions, see [REFERENCE.md](./REFERENCE.md).

- **Registry Filtering Functions**
  - [GetSubstreamsRegistry()](./REFERENCE.md#getsubstreamsregistry)
  - [GetFirehoseRegistry()](./REFERENCE.md#getfirehoseregistry)
- **Network Lookup Functions**
  - [Find(key string)](./REFERENCE.md#findkey-string)
  - [FindByGenesisBlock(blockNum uint64, blockID string)](./REFERENCE.md#findbygenesisblockblocknum-uint64-blockid-string)
  - [FindBySubstreamsEndpoint(endpoint string)](./REFERENCE.md#findbysubstreamsendpointendpoint-string)
- **Configuration Helpers**
  - [GetBytesEncoding(network *registry.Network)](./REFERENCE.md#getbytesencodingnetwork-registrynetwork)
  - [ScheduleUpdateLatestRegistry(ctx context.Context, interval time.Duration, logger *zap.Logger)](./REFERENCE.md#scheduleupdatelatestregistryctx-contextcontext-interval-timeduration-logger-zaplogger)

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

This project is licensed under the Apache License 2.0 - see the [LICENSE.md](LICENSE.md) file for details.