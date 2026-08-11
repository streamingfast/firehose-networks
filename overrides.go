package networks

import registry "github.com/pinax-network/graph-networks-libs/packages/golang/lib"

var (
	networkOverrides = []*registry.Network{
		ACMEDummyBlockchain,
	}

	serviceOverrides = []*serviceOverride{
		hoodiStreamingFast,
	}
)

// serviceOverride adds service endpoints to a network that already exists in the official
// registry, which [networkOverrides] cannot do since it only adds networks that are missing.
//
// It is meant to declare StreamingFast endpoints for networks the registry doesn't list them
// for yet. Endpoints are merged in front of the ones coming from the registry and duplicates
// are dropped, so an override turns into a no-op once the registry catches up.
type serviceOverride struct {
	// NetworkID is the [registry.Network.ID] of the network to augment, an unknown ID is ignored.
	NetworkID string

	// Firehose endpoints to add to [registry.Services.Firehose].
	Firehose []string

	// Substreams endpoints to add to [registry.Services.Substreams].
	Substreams []string
}

var (
	// StreamingFast endpoints for the Ethereum Hoodi testnet, the registry only knows about
	// the Pinax ones for now.
	hoodiStreamingFast = &serviceOverride{
		NetworkID:  "hoodi",
		Firehose:   []string{"hoodi.eth.streamingfast.io:443"},
		Substreams: []string{"hoodi.eth.streamingfast.io:443"},
	}
)

var (
	// Dummy blockchain we use for operator demonstration purposes.
	ACMEDummyBlockchain = &registry.Network{
		ID:        "acme-dummy-blockchain",
		ShortName: "Acme",
		FullName:  "Acme Dummy Blockchain",
		Aliases:   []string{"acme-dummy", "dummy-blockchain"},
		Caip2ID:   "acme:dummy-blockchain",
		Services: registry.Services{
			Firehose:   []string{"localhost:10015"},
			Substreams: []string{"localhost:10016"},
		},
		NetworkType: registry.Devnet,
		Firehose: &registry.Firehose{
			BlockType:     "sf.acme.type.v1.Block",
			BufURL:        "https://buf.build/streamingfast/firehose-acme",
			BytesEncoding: "hex",
			FirstStreamableBlock: &registry.FirstStreamableBlock{
				ID:     "0x0000000000000000000000000000000000000000000000000000000000000000",
				Height: 0,
			},
		},
	}
)
