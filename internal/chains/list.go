package chains

var Buintin = []*Chain{
	local,
	eth,
	ethSepolia,
	ethHolesky,
	arbitrum,
	arbitrumSepolia,
	opMainnet,
	opSepolia,
	bscMainnet,
	bscTestnet,
	polygon,
	polygonMubai,
	base,
	baseSepolia,
	metis,
	metisSepolia,
}

var local = &Chain{
	Name:     "local",
	Endpoint: "http://127.0.0.1:8545",
	Alias:    []string{"regtest", "devnet"},
}

var eth = &Chain{
	Name:     "eth",
	Id:       1,
	Endpoint: "https://eth.llamarpc.com",
	Explorer: "https://api.etherscan.io/api",
	Alias:    []string{"eth-mainnet", "mainnet"},
}

var ethSepolia = &Chain{
	Name:     "sepolia",
	Id:       11155111,
	Endpoint: "https://rpc.sepolia.org",
	Explorer: "https://api-sepolia.etherscan.io/api",
	Alias:    []string{"eth-sepolia"},
}

var ethHolesky = &Chain{
	Name:     "holesky",
	Id:       17000,
	Endpoint: "https://ethereum-holesky.publicnode.com",
	Explorer: "https://api-holesky.etherscan.io/api",
	Alias:    []string{"eth-holesky"},
}

var arbitrum = &Chain{
	Name:     "arbitrum",
	Id:       42161,
	Endpoint: "https://arb1.arbitrum.io/rpc",
	Explorer: "https://api.arbiscan.io/api",
	Alias:    []string{"arbi", "arbitrum-one"},
}

var arbitrumSepolia = &Chain{
	Name:     "arbitrum-sepolia",
	Id:       421614,
	Endpoint: "https://sepolia-rollup.arbitrum.io/rpc",
	Explorer: "https://api-sepolia.arbiscan.io/api",
}

// https://docs.optimism.io/chain/networks
var opMainnet = &Chain{
	Name:     "op",
	Id:       10,
	Endpoint: "https://mainnet.optimism.io",
	Explorer: "https://api-optimistic.etherscan.io/api",
	Alias:    []string{"optimistic"},
}

var opSepolia = &Chain{
	Name:     "op-sepolia",
	Id:       11155420,
	Endpoint: "https://sepolia.optimism.io",
	Explorer: "https://api-sepolia-optimism.etherscan.io/api",
}

var bscMainnet = &Chain{
	Name:     "bsc",
	Id:       56,
	Endpoint: "https://binance.llamarpc.com",
	Explorer: "https://api.bscscan.com/api",
}

var bscTestnet = &Chain{
	Name:     "bsc-testnet",
	Id:       97,
	Endpoint: "https://bsc-testnet.publicnode.com",
	Explorer: "https://api-testnet.bscscan.com/api",
}

var polygon = &Chain{
	Name:     "polygon",
	Id:       137,
	Endpoint: "https://polygon.llamarpc.com",
	Explorer: "https://api.polygonscan.com/api",
	Alias:    []string{"matic"},
}

var polygonMubai = &Chain{
	Name:     "polygon-mubai",
	Id:       80001,
	Endpoint: "https://rpc-mumbai.maticvigil.com",
	Explorer: "https://api-testnet.polygonscan.com/api",
}

// https://docs.base.org/network-information/
var base = &Chain{
	Name:     "base",
	Id:       8453,
	Endpoint: "https://mainnet.base.org",
	Explorer: "https://api.basescan.org/api",
}

var baseSepolia = &Chain{
	Name:     "base-sepolia",
	Id:       84532,
	Endpoint: "https://sepolia.base.org",
	Explorer: "https://api-sepolia.basescan.org/api",
}

var metis = &Chain{
	Name:     "metis",
	Id:       1088,
	Endpoint: "https://andromeda.metis.io",
	Explorer: "https://api.routescan.io/v2/network/mainnet/evm/1088/etherscan/api",
	Alias:    []string{"andromeda"},
}

var metisSepolia = &Chain{
	Name:     "metis-sepolia",
	Id:       59901,
	Endpoint: "https://sepolia.rpc.metisdevops.link",
	Explorer: "https://sepolia.explorer.metisdevops.link/api",
}
