package main

import (
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

func JsonrpcCall(ctx *cli.Context, rpcClient *rpc.Client) error {
	args := ctx.Args().Slice()
	jsonrpcMethod := args[0]
	jsonrpcParams, err := utils.ParseArgs(args[1:])
	if err != nil {
		return err
	}

	var result json.RawMessage
	if err := rpcClient.CallContext(ctx.Context, &result,
		jsonrpcMethod, jsonrpcParams...); err != nil && err != ethereum.NotFound {
		return err
	}
	return utils.PrintJson(os.Stdout, result, true)
}
