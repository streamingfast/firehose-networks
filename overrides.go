package networks

import registry "github.com/pinax-network/graph-networks-libs/packages/golang/lib"

var (
	networkOverrides = []*registry.Network{
		TRONMainnet,
	}
)

var (
	TRONMainnet = &registry.Network{
		ID:        "tron",
		ShortName: "Tron",
		FullName:  "Tron Mainnet",
		Aliases:   []string{"tron-mainnet"},
		Caip2ID:   "eip155:728126428",
		GraphNode: &registry.GraphNode{
			Protocol: (*registry.Protocol)(ptr("tron")),
		},
		ExplorerUrls: []string{"https://tronscan.org"},
		RPCUrls:      []string{"https://api.trongrid.io/jsonrpc", "https://tron.drpc.org", "https://rpc.ankr.com/tron_jsonrpc"},
		APIUrls: []registry.APIURL{
			{
				URL:  "https://apilist.tronscanapi.com/api/",
				Kind: "etherscan",
			},
		},
		Services: registry.Services{
			Firehose:   []string{"mainnet.tron.streamingfast.io:443"},
			Substreams: []string{"mainnet.tron.streamingfast.io:443"},
		},
		NetworkType:     "mainnet",
		IssuanceRewards: true,
		NativeToken:     ptr("TRX"),
		DocsURL:         ptr("https://developers.tron.network/"),
		Genesis: &registry.Genesis{
			Hash:   "0x00000000000000001ebf88508a03865c71d452e25f4d51194196a1d22b6653dc",
			Height: 0,
		},
		Firehose: &registry.Firehose{
			BlockType:        "sf.tron.type.v1.Block",
			EvmExtendedModel: ptr(false),
			BufURL:           "https://buf.build/streamingfast/firehose-tron",
			BytesEncoding:    "hex",
		},
	}
)
