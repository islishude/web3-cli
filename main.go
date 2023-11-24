package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/islishude/web3-cli/internal/abis"
	"github.com/islishude/web3-cli/internal/chains"
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

//go:embed brief.md
var briefDoc string

func main() {
	app := &cli.App{
		Name:        "web3-cli",
		Usage:       "web3 jsonrpc client",
		Description: briefDoc,
		Flags: []cli.Flag{
			ChainNameFlag,
			RPCFlag,
			ABINameFlag,
			ABIPathFlag,
			EthCallToFlag,
			EthCallFromFlag,
			EthCallHeightFlag,
			EthCallGasFlag,
			ExplorerApiFlag,
			ExplorerApiKeyFlag,
		},
		Commands: []*cli.Command{
			{
				Name:        "chains",
				Description: "display builtin chain config",
				Action: func(ctx *cli.Context) error {
					return utils.PrintJson(os.Stdout, chains.Buintin, true)
				},
			},
			{
				Name:        "abis",
				Description: "display builtin contracts ABI",
				Action: func(ctx *cli.Context) error {
					return utils.PrintJson(os.Stdout, abis.Builtin(), true)
				},
			},
			{
				Name:        "tools",
				Description: "handy tools",
				Subcommands: []*cli.Command{
					{
						Name:        "decode-raw-tx",
						Description: "decode raw transaction",
						Action: func(ctx *cli.Context) error {
							data, err := utils.DecodeRawTransaction(ctx.Args().First())
							if err != nil {
								return err
							}
							return utils.PrintJson(os.Stdout, json.RawMessage(data), true)
						},
					},
					{
						Name:        "new-random-address",
						Description: "create a new random address",
						Action: func(ctx *cli.Context) error {
							result, err := utils.NewRandomAddress()
							if err != nil {
								return err
							}
							return utils.PrintJson(os.Stdout, result, true)
						},
					},
				},
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Args().Len() == 0 {
				return errors.New("no jsonrpc method gived")
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			chainName := ctx.String(ChainNameFlag.Name)
			chainConf := chains.Get(chainName)
			if chainConf == nil {
				if !ctx.IsSet(RPCFlag.Name) {
					return fmt.Errorf("no rpc endpoint gived for non-builtin chain %s", chainName)
				}
				chainConf = &chains.Chain{
					Name:     chainName,
					Endpoint: ctx.String(RPCFlag.Name),
					Explorer: ctx.String(ExplorerApiFlag.Name),
				}
			} else {
				// overwrite default chain config
				if ctx.IsSet(RPCFlag.Name) {
					chainConf.Endpoint = ctx.String(RPCFlag.Name)
				}
				if ctx.IsSet(ExplorerApiFlag.Name) {
					chainConf.Explorer = ctx.String(ExplorerApiFlag.Name)
				}
			}

			rpcClient, err := rpc.DialContext(ctx.Context, chainConf.Endpoint)
			if err != nil {
				return fmt.Errorf("can not connect to node %q: %s", chainConf.Endpoint, err)
			}
			defer rpcClient.Close()

			if ctx.String(EthCallToFlag.Name) != "" {
				return ContractCall(ctx, rpcClient, chainConf)
			}
			return JsonrpcCall(ctx, rpcClient)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
