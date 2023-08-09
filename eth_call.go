package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/islishude/web3-cli/internal/abis"
	"github.com/islishude/web3-cli/internal/chains"
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

type CallMsg struct {
	From  string        `json:"from"`
	To    string        `json:"to"`
	Data  hexutil.Bytes `json:"data,omitempty"`
	Value *hexutil.Big  `json:"value,omitempty"`
	Gas   hexutil.Uint  `json:"gas,omitempty"`
}

func ContractCall(ctx *cli.Context, rpcClient *rpc.Client, chain *chains.Chain) (err error) {
	callMsg := CallMsg{
		From: ctx.String(EthCallFromFlag.Name),
		To:   ctx.String(EthCallToFlag.Name),
	}

	if ctx.IsSet(EthCallGasFlag.Name) {
		callMsg.Gas = hexutil.Uint(ctx.Uint64(EthCallGasFlag.Name))
	}

	if ctx.IsSet(EthCallValueFlag.Name) {
		callMsg.Value = (*hexutil.Big)(hexutil.MustDecodeBig(ctx.String(EthCallValueFlag.Name)))
	}

	abiIns, err := getABI(ctx, chain, callMsg.To)
	if err != nil {
		return err
	}

	callMsg.Data, err = abis.Pack(abiIns, ctx.Args().Slice())
	if err != nil {
		return err
	}

	jsonrpcParam := []any{callMsg, ctx.String(EthCallHeightFlag.Name)}

	var output hexutil.Bytes
	if err := rpcClient.CallContext(ctx.Context, &output, "eth_call", jsonrpcParam...); err != nil {
		return err
	}

	methodIns := abiIns.Methods[ctx.Args().First()]
	if len(methodIns.Outputs) > 0 && len(output) == 0 {
		jsonrpcParam[0] = callMsg.To
		if err = rpcClient.CallContext(ctx.Context, &output, "eth_getCode", jsonrpcParam...); err != nil {
			return err
		}

		if len(output) == 0 {
			return fmt.Errorf("gived address %s is not a contract", callMsg.To)
		}
	}

	result, err := abis.Unpack(abiIns, ctx.Args().First(), output)
	if err != nil {
		return err
	}

	return utils.PrintJson(os.Stdout, result, true)
}

func getABI(ctx *cli.Context, chain *chains.Chain, contAddr string) (*abi.ABI, error) {
	abiIns := abis.Get(ctx.String(ABINameFlag.Name))
	if abiIns != nil {
		return abiIns, nil
	}

	// you can always overwrite builtin ABI
	fetchURL, isExplorer := ctx.String(ABIPathFlag.Name), false
	if fetchURL == "" {
		expURL, err := utils.URLToGetABI(chain.Explorer, contAddr, ctx.String(ExplorerApiKeyFlag.Name))
		if err != nil {
			return nil, err
		}
		fetchURL, isExplorer = expURL, true
	}

	return abis.Fetch(fetchURL, isExplorer)
}
