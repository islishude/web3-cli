package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/rpc"
)

const helpText = "web3-cli - web3 jsonrpc client\n\nUsage: web3-cli method [param...]\n\nDefault web3 server endpoint is `http://localhost:8545`, you can set `web3` env value to change it."

func main() {
	var method string
	var params []interface{}

	switch len(os.Args) {
	case 1:
		log.Println(helpText)
		return
	case 2:
		method = os.Args[1]
		if method == "-h" || method == "--help" {
			log.Println(helpText)
			return
		}
	default:
		method = os.Args[1]
		var err error
		params, err = parseArgs(os.Args[2:])
		if err != nil {
			log.Fatalln(err)
		}
	}

	endpoint := os.Getenv("web3")
	if endpoint == "" {
		endpoint = "http://127.0.0.1:8545"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ethclient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		log.Fatalf("can't connect to web3 server %q: %s\n", endpoint, err)
		return
	}
	defer ethclient.Close()

	var result json.RawMessage
	if err := ethclient.CallContext(ctx, &result, method, params...); err != nil && err != ethereum.NotFound {
		log.Fatalln(err)
		return
	}

	out := json.NewEncoder(os.Stdout)
	out.SetIndent("", "  ")
	_ = out.Encode(result)
}
