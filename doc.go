// Package networks provides a wrapper around The Graph's Networks Registry library
// specifically designed for Substreams and Firehose ecosystem products.
//
// This package extends the base Networks Registry functionality with:
//   - Helper utilities commonly used across Substreams/Firehose projects
//   - Support for custom network definitions for development and testing
//   - Automatic fallback to embedded registry data when remote registry is unavailable
//   - Enhanced network lookup and filtering capabilities
//
// # Usage
//
// Get all networks with Substreams endpoints:
//
//	substreamsNetworks := networks.GetSubstreamsRegistry()
//	for id, network := range substreamsNetworks {
//	    fmt.Printf("Network %s has Substreams endpoints: %v\n", id, network.Services.Substreams)
//	}
//
// Get all networks with Firehose endpoints:
//
//	firehoseNetworks := networks.GetFirehoseRegistry()
//	for id, network := range firehoseNetworks {
//	    fmt.Printf("Network %s has Firehose endpoints: %v\n", id, network.Services.Firehose)
//	}
//
// Find a network by various identifiers:
//
//	// By network ID
//	network := networks.Find("ethereum-mainnet")
//
//	// By alias
//	network = networks.Find("eth")
//
//	// By full name
//	network = networks.Find("Ethereum Mainnet")
//
// # Fallback Mechanism
//
// The package automatically handles network registry availability. If the remote
// registry cannot be loaded, it falls back to an embedded JSON file and launches
// a background process to retry loading the latest registry with exponential backoff.
//
// # Custom Networks
//
// The package supports custom network overrides for development and testing purposes.
// These are automatically merged with the official registry data.
package networks
