package main

import (
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

const (
	CallCategory = "ETH CALL"
)

var (
	ChainNameFlag = &cli.StringFlag{
		Name:    "chain",
		Value:   "local",
		Aliases: []string{"c"},
		Usage:   "chain id or chain name, run `web3-cli chains` to get built-in chains",
		EnvVars: []string{"web3"},
	}

	RPCFlag = &cli.StringFlag{
		Name:    "rpc",
		Aliases: []string{"r"},
		Usage:   "custom rpc endpoint to connect",
	}

	NoTranformFlag = &cli.BoolFlag{
		Name:    "no-tranform",
		Aliases: []string{"n"},
		Usage:   "Use eth json rpc mode(e.g. tranform number to hex)",
	}

	ABINameFlag = &cli.StringFlag{
		Name:     "abi-name",
		Usage:    "if you want to do a contract call, you should provide this",
		Category: CallCategory,
	}

	ABIPathFlag = &cli.StringFlag{
		Name:     "abi-path",
		Usage:    "use your local abi path, or a url from a http server",
		Category: CallCategory,
	}

	EthCallToFlag = &cli.StringFlag{
		Name:     "call-to",
		Aliases:  []string{"contract", "ct"},
		Usage:    "contract address to read or write",
		Action:   utils.EthCallToFlagAction,
		Category: CallCategory,
	}

	EthCallFromFlag = &cli.StringFlag{
		Name:     "call-from",
		Usage:    "from address for contract call",
		Action:   utils.EthCallFromFlagAction,
		Category: CallCategory,
	}

	EthCallHeightFlag = &cli.StringFlag{
		Name:     "call-height",
		Value:    "latest",
		Usage:    "height for contract call",
		Action:   utils.EthCallHeightFlagAction,
		Category: CallCategory,
	}

	EthCallGasFlag = &cli.StringFlag{
		Name:     "call-gas",
		Value:    "0",
		Usage:    "gas limit for contract call",
		Action:   utils.EthCallGasFlagAction,
		Category: CallCategory,
	}

	EthCallValueFlag = &cli.StringFlag{
		Name:     "call-value",
		Value:    "0",
		Usage:    "ether value for contract call",
		Action:   utils.EthCallValueFlagAction,
		Category: CallCategory,
	}

	ExplorerApiFlag = &cli.StringFlag{
		Name:     "explorer-api",
		Usage:    "explorer api, it's used for fetching ABI",
		Category: CallCategory,
	}

	ExplorerApiKeyFlag = &cli.StringFlag{
		Name:    "explorer-api-key",
		Usage:   "explorer api key",
		EnvVars: []string{"WEB3_CLI_EXPLORER_API_KEY"},
	}

	NewAddressPrefixFlag = &cli.StringFlag{
		Name:   "prefix",
		Usage:  "address prefix",
		Action: utils.HexStringValidAction,
	}

	NewAddressSuffixFlag = &cli.StringFlag{
		Name:   "suffix",
		Usage:  "address suffix",
		Action: utils.HexStringValidAction,
	}
)
