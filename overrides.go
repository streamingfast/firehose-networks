package networks

import registry "github.com/pinax-network/graph-networks-libs/packages/golang/lib"

var (
	networkOverrides = []*registry.Network{
		ACMEDummyBlockchain,
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
