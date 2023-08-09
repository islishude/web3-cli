package main

import (
	"fmt"

	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

var (
	ChainNameFlag = &cli.StringFlag{
		Name:    "chain",
		Value:   "local",
		Aliases: []string{"c"},
		Usage:   "chain id or chain name, run `web3-cli chains` to get built-in chain configuration",
		EnvVars: []string{"WEB3_CLI_CHAIN", "web3" /*compatible*/},
	}

	RPCFlag = &cli.StringFlag{
		Name:    "rpc",
		Aliases: []string{"r"},
		Usage:   "custom rpc endpoint to connect",
		EnvVars: []string{"WEB3_CLI_RPC"},
	}

	ABINameFlag = &cli.StringFlag{
		Name:    "abi-name",
		Aliases: []string{"an"},
		Usage:   "if you want to do a contract call, you should provide this",
		EnvVars: []string{"WEB3_CLI_ABI_NAME"},
	}

	ABIPathFlag = &cli.StringFlag{
		Name:    "abi-path",
		Aliases: []string{"ap"},
		Usage:   "use your local abi path, or a url from a http server",
		EnvVars: []string{"WEB3_CLI_ABI_PATH"},
	}

	EthCallToFlag = &cli.StringFlag{
		Name:    "call-to",
		Aliases: []string{"contract", "ct"},
		Usage:   "contract address to read or write",
		EnvVars: []string{"WEB3_CLI_CONTRACT_ADDRESS"},
		Action:  utils.EthCallToFlagAction,
	}

	EthCallFromFlag = &cli.StringFlag{
		Name:    "call-from",
		Aliases: []string{"cf"},
		Value:   "0x0000000000000000000000000000000000000000",
		Usage:   "from address for contract call",
		EnvVars: []string{"WEB3_CLI_CONTRACT_FROM_ADDRESS"},
		Action:  utils.EthCallFromFlagAction,
	}

	EthCallHeightFlag = &cli.StringFlag{
		Name:    "call-height",
		Aliases: []string{"ch"},
		Value:   "latest",
		Usage:   "height for contract call",
		EnvVars: []string{"WEB3_CLI_CONTRACT_FROM_ADDRESS"},
		Action:  utils.EthCallHeightFlagAction,
	}

	EthCallGasFlag = &cli.StringFlag{
		Name:    "call-gas",
		Aliases: []string{"cg"},
		Value:   "0",
		Usage:   "gas limit for contract call",
		EnvVars: []string{"WEB3_CLI_CONTRACT_CALL_GAS"},
		Action: func(ctx *cli.Context, s string) error {
			if v := utils.ToBigInt(s); v != nil {
				return ctx.Set("call-gas", v.String())
			}
			return fmt.Errorf("invalid number: %s", s)
		},
	}

	EthCallValueFlag = &cli.StringFlag{
		Name:    "call-value",
		Aliases: []string{"cg"},
		Value:   "0",
		Usage:   "ether value for contract call",
		EnvVars: []string{"WEB3_CLI_CONTRACT_CALL_VALUE"},
		Action:  utils.EthCallValueFlagAction,
	}

	ExplorerApiFlag = &cli.StringFlag{
		Name:    "explorer-api",
		Aliases: []string{"ea"},
		Usage:   "explorer api, it's used for fetching ABI",
		EnvVars: []string{"WEB3_CLI_EXPLORER_API"},
	}

	ExplorerApiKeyFlag = &cli.StringFlag{
		Name:    "explorer-api-key",
		Aliases: []string{"ek"},
		Usage:   "explorer api key",
		EnvVars: []string{"WEB3_CLI_EXPLORER_API_KEY"},
	}
)
