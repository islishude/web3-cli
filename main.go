package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

const helpText = "web3-cli - web3 jsonrpc client tools\n\nUsage: web3-cli method [param...]\n\nDefault web3 server endpoint is `http://localhost:8545`, you can set `web3` env value to change it."

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
		for _, p := range os.Args[2:] {
			switch {
			case regexp.MustCompile(`^[0-9]+$`).MatchString(p):
				if v, ok := new(big.Int).SetString(p, 10); ok {
					params = append(params, "0x"+v.Text(16))
				} else {
					params = append(params, p)
				}
			case p == "true":
				params = append(params, true)
			case p == "false":
				params = append(params, false)
			default:
				params = append(params, p)
			}
		}
	}

	endpoint := os.Getenv("web3")
	if endpoint == "" {
		log.Println("Using default localhost rpc endpoint")
		endpoint = "http://127.0.0.1:8545"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ethclient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		log.Printf("can't connect web3 server %q: %s\n", endpoint, err)
		return
	}
	defer ethclient.Close()

	var result json.RawMessage
	if err := ethclient.CallContext(ctx, &result, method, params...); err != nil {
		log.Printf("Call %s failed with params %v: %s\n", method, params, err)
		return
	}

	out := json.NewEncoder(os.Stdout)
	out.SetIndent("", "  ")
	_ = out.Encode(result)
}
