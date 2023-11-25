package main

import (
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/islishude/web3-cli/internal/abis"
	"github.com/islishude/web3-cli/internal/chains"
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

type EthCallMsg struct {
	From  string        `json:"from,omitempty"`
	To    string        `json:"to"`
	Data  hexutil.Bytes `json:"data,omitempty"`
	Value *hexutil.Big  `json:"value,omitempty"`
	Gas   hexutil.Uint  `json:"gas,omitempty"`
}

func ContractCall(ctx *cli.Context, rpcClient *rpc.Client, chain *chains.Chain, logger io.Writer) (err error) {
	callMsg := EthCallMsg{To: ctx.String(EthCallToFlag.Name)}

	if v := ctx.String(EthCallFromFlag.Name); v != "'" {
		callMsg.From = v
	}

	if v := ctx.Uint64(EthCallGasFlag.Name); v != 0 {
		callMsg.Gas = hexutil.Uint(v)
	}

	if v := ctx.String(EthCallValueFlag.Name); v != "" {
		callMsg.Value = (*hexutil.Big)(hexutil.MustDecodeBig(v))
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

	return utils.PrintJson(logger, result, true)
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
