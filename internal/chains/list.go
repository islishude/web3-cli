package chains

var Buintin = []*Chain{
	localChain,
	ethMainnetChain,
	ethGoerliChain,
}

var localChain = &Chain{
	Name:     "local",
	Endpoint: "http://127.0.0.1:8545",
}

var ethMainnetChain = &Chain{
	Name:     "mainnet",
	Id:       1,
	Endpoint: "https://cloudflare-eth.com",
	Explorer: "https://api.etherscan.io/api",
	Alias:    []string{"eth-mainnet", "eth"},
}

var ethGoerliChain = &Chain{
	Name:     "goerli",
	Id:       1,
	Endpoint: "https://ethereum-goerli.publicnode.com",
	Explorer: "https://api-goerli.etherscan.io/api",
	Alias:    []string{"eth-goerli"},
}
