package networks

import (
	"regexp"
	"testing"

	registry "github.com/pinax-network/graph-networks-libs/packages/golang/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkRegistry_Find(t *testing.T) {
	net1 := &registry.Network{ID: "mainnet", ShortName: "ETH", FullName: "Ethereum Mainnet", Aliases: []string{"eth", "ethereum"}}
	net2 := &registry.Network{ID: "arbitrum", ShortName: "ARB", FullName: "Arbitrum One", Aliases: []string{"arb", "arbitrum-one"}}
	net3 := &registry.Network{ID: "custom", ShortName: "MYC", FullName: "My Custom Chain", Aliases: []string{"mychain"}}

	r := NetworkRegistry{
		"mainnet":  net1,
		"arbitrum": net2,
		"custom":   net3,
	}

	t.Run("find by id", func(t *testing.T) {
		assert.Equal(t, net1, r.Find("mainnet"))
		assert.Equal(t, net2, r.Find("arbitrum"))
	})
	t.Run("find by alias", func(t *testing.T) {
		assert.Equal(t, net1, r.Find("eth"))
		assert.Equal(t, net1, r.Find("ethereum"))
		assert.Equal(t, net2, r.Find("arb"))
		assert.Equal(t, net2, r.Find("arbitrum-one"))
		assert.Equal(t, net3, r.Find("mychain"))
	})
	t.Run("find by FullName", func(t *testing.T) {
		assert.Equal(t, net1, r.Find("Ethereum Mainnet"))
		assert.Equal(t, net2, r.Find("Arbitrum One"))
		assert.Equal(t, net3, r.Find("My Custom Chain"))
	})
	t.Run("find by ShortName", func(t *testing.T) {
		assert.Equal(t, net1, r.Find("ETH"))
		assert.Equal(t, net2, r.Find("ARB"))
		assert.Equal(t, net3, r.Find("MYC"))
	})
	t.Run("not found", func(t *testing.T) {
		assert.Nil(t, r.Find("notfound"))
	})
}

func TestNetworkRegistry_FindAll(t *testing.T) {
	net1 := &registry.Network{ID: "mainnet", ShortName: "ETH", FullName: "Ethereum Mainnet", Aliases: []string{"eth", "ethereum"}}
	net2 := &registry.Network{ID: "arbitrum", ShortName: "ARB", FullName: "Arbitrum One", Aliases: []string{"arb", "arbitrum-one"}}
	net3 := &registry.Network{ID: "custom", ShortName: "ETH", FullName: "My Custom Chain", Aliases: []string{"mychain"}}

	r := NetworkRegistry{
		"mainnet":  net1,
		"arbitrum": net2,
		"custom":   net3,
	}

	t.Run("find all by id", func(t *testing.T) {
		results := r.FindAll("mainnet")
		assert.Len(t, results, 1)
		assert.Equal(t, net1, results[0])
	})
	t.Run("find all by alias", func(t *testing.T) {
		results := r.FindAll("eth")
		assert.Len(t, results, 1)
		assert.Equal(t, net1, results[0])

		results = r.FindAll("arbitrum-one")
		assert.Len(t, results, 1)
		assert.Equal(t, net2, results[0])
	})
	t.Run("find all by FullName", func(t *testing.T) {
		results := r.FindAll("Ethereum Mainnet")
		assert.Len(t, results, 1)
		assert.Equal(t, net1, results[0])
	})
	t.Run("find all by ShortName with multiple matches", func(t *testing.T) {
		results := r.FindAll("ETH")
		assert.Len(t, results, 2)
		assert.Contains(t, results, net1)
		assert.Contains(t, results, net3)
	})
	t.Run("not found", func(t *testing.T) {
		results := r.FindAll("notfound")
		assert.Empty(t, results)
	})
	t.Run("results are sorted by ID", func(t *testing.T) {
		results := r.FindAll("ETH")
		assert.Len(t, results, 2)
		assert.Equal(t, "custom", results[0].ID)
		assert.Equal(t, "mainnet", results[1].ID)
	})
}

func TestNetworkRegistry_Search(t *testing.T) {
	net1 := &registry.Network{ID: "mainnet", ShortName: "ETH", FullName: "Ethereum Mainnet", Aliases: []string{"eth", "ethereum"}}
	net2 := &registry.Network{ID: "arbitrum", ShortName: "ARB", FullName: "Arbitrum One", Aliases: []string{"arb", "arbitrum-one"}}
	net3 := &registry.Network{ID: "optimism", ShortName: "OPT", FullName: "Optimism Mainnet", Aliases: []string{"opt", "optimism-mainnet"}}
	net4 := &registry.Network{ID: "polygon", ShortName: "POL", FullName: "Polygon Mainnet", Aliases: []string{"matic"}}

	r := NetworkRegistry{
		"mainnet":  net1,
		"arbitrum": net2,
		"optimism": net3,
		"polygon":  net4,
	}

	t.Run("search by ID pattern", func(t *testing.T) {
		re := regexp.MustCompile("^main")
		results := r.Search(re)
		assert.Len(t, results, 1)
		assert.Equal(t, net1, results[0])
	})
	t.Run("search by ShortName pattern", func(t *testing.T) {
		re := regexp.MustCompile("^AR")
		results := r.Search(re)
		assert.Len(t, results, 1)
		assert.Equal(t, net2, results[0])
	})
	t.Run("search by FullName pattern", func(t *testing.T) {
		re := regexp.MustCompile("Mainnet$")
		results := r.Search(re)
		assert.Len(t, results, 3)
		assert.Contains(t, results, net1)
		assert.Contains(t, results, net3)
		assert.Contains(t, results, net4)
	})
	t.Run("search by alias pattern", func(t *testing.T) {
		re := regexp.MustCompile("^arb")
		results := r.Search(re)
		assert.Len(t, results, 1)
		assert.Equal(t, net2, results[0])
	})
	t.Run("search with pattern matching multiple fields", func(t *testing.T) {
		re := regexp.MustCompile("(?i)eth")
		results := r.Search(re)
		assert.Len(t, results, 1)
		assert.Equal(t, net1, results[0])
	})
	t.Run("no matches", func(t *testing.T) {
		re := regexp.MustCompile("^notfound")
		results := r.Search(re)
		assert.Empty(t, results)
	})
	t.Run("results are sorted by ID", func(t *testing.T) {
		re := regexp.MustCompile(".*")
		results := r.Search(re)
		assert.Len(t, results, 4)
		assert.Equal(t, "arbitrum", results[0].ID)
		assert.Equal(t, "mainnet", results[1].ID)
		assert.Equal(t, "optimism", results[2].ID)
		assert.Equal(t, "polygon", results[3].ID)
	})
}

func TestAllLegacyChainConfigKeysPresent(t *testing.T) {
	legacyKeys := []string{
		"mainnet", "bnb", "polygon", "amoy", "arbitrum", "holesky", "sepolia", "optimism", "avalanche", "chapel",
		"injective-mainnet", "injective-testnet", "starknet-mainnet", "starknet-testnet", "solana-mainnet-beta",
		"mantra-testnet", "mantra-mainnet", "stellar-testnet", "stellar", "sei-mainnet",
	}

	for _, key := range legacyKeys {
		net := Find(key)
		assert.NotNilf(t, net, "Network with key %q should be present in GetSubstreamsRegistry()", key)
	}
}

func TestGetSubstreamsRegistry(t *testing.T) {
	networks := GetSubstreamsRegistry()
	assert.NotEmpty(t, networks, "Should return at least one network with Substreams endpoint")
	for id, net := range networks {
		assert.Greater(t, len(net.Services.Substreams), 0, "Network %q should have at least one Substreams endpoint", id)
	}
	// Known networks with Substreams endpoints (should be present)
	for _, key := range []string{"mainnet", "optimism", "arbitrum", "polygon", "bnb", "avalanche"} {
		assert.NotNilf(t, networks.Find(key), "Network %q should be present in Substreams registry", key)
	}
	// Known networks without Substreams endpoints (should NOT be present)
	for _, key := range []string{"cronos", "clover", "aurora", "celo"} {
		assert.Nilf(t, networks.Find(key), "Network %q should NOT be present in Substreams registry", key)
	}
}

func TestGetFirehoseRegistry(t *testing.T) {
	networks := GetFirehoseRegistry()
	assert.NotEmpty(t, networks, "Should return at least one network with Firehose endpoint")
	for id, net := range networks {
		assert.Greater(t, len(net.Services.Firehose), 0, "Network %q should have at least one Firehose endpoint", id)
	}
	// Known networks with Firehose endpoints (should be present)
	for _, key := range []string{"mainnet", "optimism", "arbitrum", "polygon", "bnb", "avalanche"} {
		assert.NotNilf(t, networks.Find(key), "Network %q should be present in Firehose registry", key)
	}
	// Known networks without Firehose endpoints (should NOT be present)
	for _, key := range []string{"cronos", "clover", "aurora", "celo"} {
		assert.Nilf(t, networks.Find(key), "Network %q should NOT be present in Firehose registry", key)
	}
}

func TestNetworkRegistry_FindByGenesisBlock(t *testing.T) {
	networks := getRegistryNetworksFull()
	const moonbeamID = "moonbeam"
	const moonbeamGenesisHash = "0x7e6b3bbed86828a558271c9c9f62354b1d8b5aa15ff85fd6f1e7cbe9af9dde7e"
	const moonbeamGenesisHeight = 0

	net := networks.FindByGenesisBlock(moonbeamGenesisHeight, moonbeamGenesisHash)
	assert.NotNil(t, net, "Should find Moonbeam by genesis block")
	assert.Equal(t, moonbeamID, net.ID)

	// Not found case
	notFound := networks.FindByGenesisBlock(12345, "0xdeadbeef")
	assert.Nil(t, notFound, "Should return nil for unknown genesis block")
}

func TestGetBytesEncoding(t *testing.T) {
	networks := getRegistryNetworksFull()

	t.Run("returns correct encoding for mainnet", func(t *testing.T) {
		net := networks.Find("mainnet")
		assert.NotNil(t, net)
		assert.Equal(t, registry.Hex, GetBytesEncoding(net))
	})

	t.Run("returns correct encoding for optimism", func(t *testing.T) {
		net := networks.Find("optimism")
		assert.NotNil(t, net)
		assert.Equal(t, registry.Hex, GetBytesEncoding(net))
	})

	t.Run("returns Hex for nil network", func(t *testing.T) {
		assert.Equal(t, registry.Hex, GetBytesEncoding(nil))
	})

	t.Run("returns Hex for network without Firehose", func(t *testing.T) {
		net := &registry.Network{ID: "no-firehose"}
		assert.Equal(t, registry.Hex, GetBytesEncoding(net))
	})
}

func TestFindBySubstreamsEndpoint(t *testing.T) {
	substreamsRegistry := GetSubstreamsRegistry()

	t.Run("finds mainnet by endpoint", func(t *testing.T) {
		mainnetEndpoints := []string{
			"eth.substreams.pinax.network:443",
			"mainnet.eth.streamingfast.io:443",
		}
		for _, ep := range mainnetEndpoints {
			net := substreamsRegistry.FindBySubstreamsEndpoint(ep)
			assert.NotNilf(t, net, "Should find mainnet for endpoint %q", ep)
			assert.Equal(t, "mainnet", net.ID)
		}
	})

	t.Run("finds optimism by endpoint", func(t *testing.T) {
		optimismEndpoints := []string{
			"mainnet.optimism.streamingfast.io:443",
			"optimism.substreams.pinax.network:443",
		}
		for _, ep := range optimismEndpoints {
			net := substreamsRegistry.FindBySubstreamsEndpoint(ep)
			assert.NotNilf(t, net, "Should find optimism for endpoint %q", ep)
			assert.Equal(t, "optimism", net.ID)
		}
	})

	t.Run("returns nil for unknown endpoint", func(t *testing.T) {
		net := substreamsRegistry.FindBySubstreamsEndpoint("unknown.endpoint:1234")
		assert.Nil(t, net)
	})

	t.Run("returns nil for empty endpoint", func(t *testing.T) {
		net := substreamsRegistry.FindBySubstreamsEndpoint("")
		assert.Nil(t, net)
	})

	t.Run("returns nil for network with no Substreams endpoints", func(t *testing.T) {
		// Add a network with no Substreams endpoints
		net := &registry.Network{ID: "no-substreams", Services: registry.Services{Substreams: []string{}}}
		r := NetworkRegistry{"no-substreams": net}
		assert.Nil(t, r.FindBySubstreamsEndpoint("any.endpoint:443"))
	})
}

func TestGetSubstreamsEndpoint(t *testing.T) {
	t.Run("prioritizes streamingfast.io endpoint", func(t *testing.T) {
		endpoint := GetSubstreamsEndpoint("mainnet")
		assert.NotEmpty(t, endpoint, "Should return an endpoint for mainnet")
		assert.Contains(t, endpoint, "streamingfast.io", "Should prioritize streamingfast.io endpoint")
		assert.Equal(t, "mainnet.eth.streamingfast.io:443", endpoint)
	})

	t.Run("returns first endpoint when no streamingfast.io endpoint", func(t *testing.T) {
		// Create a test network with no streamingfast.io endpoint
		net := &registry.Network{
			ID: "test-no-sf",
			Services: registry.Services{
				Substreams: []string{
					"test.pinax.network:443",
					"test.other.network:443",
				},
			},
		}

		// Temporarily add to registry
		reg := getRegistryNetworksFull()
		reg["test-no-sf"] = net
		defer delete(reg, "test-no-sf")

		endpoint := GetSubstreamsEndpoint("test-no-sf")
		assert.Equal(t, "test.pinax.network:443", endpoint)
	})

	t.Run("returns empty string for unknown network", func(t *testing.T) {
		endpoint := GetSubstreamsEndpoint("unknown-network")
		assert.Empty(t, endpoint)
	})

	t.Run("returns empty string for network with no substreams endpoints", func(t *testing.T) {
		// Create a test network with no substreams endpoints
		net := &registry.Network{
			ID: "test-no-substreams",
			Services: registry.Services{
				Substreams: []string{},
			},
		}

		// Temporarily add to registry
		reg := getRegistryNetworksFull()
		reg["test-no-substreams"] = net
		defer delete(reg, "test-no-substreams")

		endpoint := GetSubstreamsEndpoint("test-no-substreams")
		assert.Empty(t, endpoint)
	})
}

func TestGetFirehoseEndpoint(t *testing.T) {
	t.Run("prioritizes streamingfast.io endpoint", func(t *testing.T) {
		endpoint := GetFirehoseEndpoint("mainnet")
		assert.NotEmpty(t, endpoint, "Should return an endpoint for mainnet")
		assert.Contains(t, endpoint, "streamingfast.io", "Should prioritize streamingfast.io endpoint")
		assert.Equal(t, "mainnet.eth.streamingfast.io:443", endpoint)
	})

	t.Run("returns first endpoint when no streamingfast.io endpoint", func(t *testing.T) {
		// Create a test network with no streamingfast.io endpoint
		net := &registry.Network{
			ID: "test-no-sf",
			Services: registry.Services{
				Firehose: []string{
					"test.pinax.network:443",
					"test.other.network:443",
				},
			},
		}

		// Temporarily add to registry
		reg := getRegistryNetworksFull()
		reg["test-no-sf"] = net
		defer delete(reg, "test-no-sf")

		endpoint := GetFirehoseEndpoint("test-no-sf")
		assert.Equal(t, "test.pinax.network:443", endpoint)
	})

	t.Run("returns empty string for unknown network", func(t *testing.T) {
		endpoint := GetFirehoseEndpoint("unknown-network")
		assert.Empty(t, endpoint)
	})

	t.Run("returns empty string for network with no firehose endpoints", func(t *testing.T) {
		// Create a test network with no firehose endpoints
		net := &registry.Network{
			ID: "test-no-firehose",
			Services: registry.Services{
				Firehose: []string{},
			},
		}

		// Temporarily add to registry
		reg := getRegistryNetworksFull()
		reg["test-no-firehose"] = net
		defer delete(reg, "test-no-firehose")

		endpoint := GetFirehoseEndpoint("test-no-firehose")
		assert.Empty(t, endpoint)
	})
}

func TestNetworkRegistry_addServiceEndpoints(t *testing.T) {
	newRegistry := func() NetworkRegistry {
		return NetworkRegistry{
			"hoodi": &registry.Network{
				ID: "hoodi",
				Services: registry.Services{
					Firehose:   []string{"hoodi.firehose.pinax.network:443"},
					Substreams: []string{"hoodi.substreams.pinax.network:443"},
				},
			},
		}
	}

	t.Run("prepends endpoints to an existing network", func(t *testing.T) {
		r := newRegistry()
		r.addServiceEndpoints(&serviceOverride{
			NetworkID:  "hoodi",
			Firehose:   []string{"hoodi.eth.streamingfast.io:443"},
			Substreams: []string{"hoodi.eth.streamingfast.io:443"},
		})

		assert.Equal(t, []string{"hoodi.eth.streamingfast.io:443", "hoodi.firehose.pinax.network:443"}, r["hoodi"].Services.Firehose)
		assert.Equal(t, []string{"hoodi.eth.streamingfast.io:443", "hoodi.substreams.pinax.network:443"}, r["hoodi"].Services.Substreams)
	})

	t.Run("does not duplicate endpoints already in the registry", func(t *testing.T) {
		r := newRegistry()
		override := &serviceOverride{
			NetworkID:  "hoodi",
			Firehose:   []string{"hoodi.firehose.pinax.network:443"},
			Substreams: []string{"hoodi.substreams.pinax.network:443"},
		}

		r.addServiceEndpoints(override)
		r.addServiceEndpoints(override)

		assert.Equal(t, []string{"hoodi.firehose.pinax.network:443"}, r["hoodi"].Services.Firehose)
		assert.Equal(t, []string{"hoodi.substreams.pinax.network:443"}, r["hoodi"].Services.Substreams)
	})

	t.Run("leaves untouched services the override does not define", func(t *testing.T) {
		r := newRegistry()
		r.addServiceEndpoints(&serviceOverride{
			NetworkID: "hoodi",
			Firehose:  []string{"hoodi.eth.streamingfast.io:443"},
		})

		assert.Equal(t, []string{"hoodi.eth.streamingfast.io:443", "hoodi.firehose.pinax.network:443"}, r["hoodi"].Services.Firehose)
		assert.Equal(t, []string{"hoodi.substreams.pinax.network:443"}, r["hoodi"].Services.Substreams)
	})

	t.Run("ignores unknown network and invalid input", func(t *testing.T) {
		r := newRegistry()
		r.addServiceEndpoints(&serviceOverride{NetworkID: "unknown-network", Firehose: []string{"unknown.streamingfast.io:443"}})
		r.addServiceEndpoints(&serviceOverride{Firehose: []string{"unknown.streamingfast.io:443"}})
		r.addServiceEndpoints(nil)

		assert.Equal(t, newRegistry(), r)
	})
}

func TestServiceOverrides_Hoodi(t *testing.T) {
	// Loaded from the embedded JSON so the assertions are not affected by the live registry
	reg, err := loadRegistry(fromEmbeddedJSON)
	require.NoError(t, err)

	net := reg.Find("hoodi")
	require.NotNil(t, net, "Network %q should be present in the registry", "hoodi")

	assert.Equal(t, []string{"hoodi.eth.streamingfast.io:443", "hoodi.firehose.pinax.network:443"}, net.Services.Firehose)
	assert.Equal(t, []string{"hoodi.eth.streamingfast.io:443", "hoodi.substreams.pinax.network:443"}, net.Services.Substreams)
}

func TestMergeEndpoints(t *testing.T) {
	preferred := []string{"a", "b"}
	existing := []string{"b", "c"}

	assert.Equal(t, []string{"a", "b", "c"}, mergeEndpoints(preferred, existing))
	assert.Equal(t, []string{"a", "b"}, preferred, "Input slice must not be mutated")
	assert.Equal(t, []string{"b", "c"}, existing, "Input slice must not be mutated")

	assert.Equal(t, []string{"b", "c"}, mergeEndpoints(nil, existing))
	assert.Equal(t, []string{"a", "b"}, mergeEndpoints(preferred, nil))
	assert.Nil(t, mergeEndpoints(nil, nil))
}
