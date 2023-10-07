package chains

var Buintin = []*Chain{
	local,
	eth,
	ethGoerli,
	ethSepolia,
	ethHolesky,
	arbitrum,
	arbitrumGoerli,
	opMainnet,
	opGoerli,
	bscMainnet,
	bscTestnet,
	polygon,
	polygonMubai,
	base,
	baseGoerli,
	metis,
}

var local = &Chain{
	Name:     "local",
	Endpoint: "http://127.0.0.1:8545",
	Alias:    []string{"regtest", "devnet"},
}

var eth = &Chain{
	Name:     "eth",
	Id:       1,
	Endpoint: "https://cloudflare-eth.com",
	Explorer: "https://api.etherscan.io/api",
	Alias:    []string{"eth-mainnet", "mainnet"},
}

var ethGoerli = &Chain{
	Name:     "goerli",
	Id:       5,
	Endpoint: "https://ethereum-goerli.publicnode.com",
	Explorer: "https://api-goerli.etherscan.io/api",
	Alias:    []string{"eth-goerli"},
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
	Explorer: "https://holesky.etherscan.io/api",
	Alias:    []string{"eth-holesky"},
}

var arbitrum = &Chain{
	Name:     "arbitrum",
	Id:       42161,
	Endpoint: "https://arb1.arbitrum.io/rpc",
	Explorer: "https://api.arbiscan.io/api",
	Alias:    []string{"arbi", "arbitrum-one"},
}

var arbitrumGoerli = &Chain{
	Name:     "arbitrum-goerli",
	Id:       421613,
	Endpoint: "https://goerli-rollup.arbitrum.io/rpc",
	Explorer: "https://api-goerli.arbiscan.io/api",
	Alias:    []string{"arbi-goerli"},
}

var opMainnet = &Chain{
	Name:     "op",
	Id:       10,
	Endpoint: "https://mainnet.optimism.io",
	Explorer: "https://api-optimistic.etherscan.io/api",
	Alias:    []string{"optimistic"},
}

var opGoerli = &Chain{
	Name:     "op-goerli",
	Id:       420,
	Endpoint: "https://goerli.optimism.io",
	Explorer: "https://api-goerli-optimism.etherscan.io/api",
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
	Endpoint: "https://data-seed-prebsc-1-s1.bnbchain.org:8545",
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

var base = &Chain{
	Name:     "base",
	Id:       8453,
	Endpoint: "https://mainnet.base.org",
	Explorer: "https://api.basescan.org/api",
}

var baseGoerli = &Chain{
	Name:     "base-goerli",
	Id:       84531,
	Endpoint: "https://goerli.base.org",
	Explorer: "https://api-goerli.basescan.org/api",
}

var metis = &Chain{
	Name:     "metis",
	Id:       1088,
	Endpoint: "https://andromeda.metis.io",
	Explorer: "https://api.routescan.io/v2/network/mainnet/evm/1088/etherscan/api",
	Alias:    []string{"andromeda"},
}
