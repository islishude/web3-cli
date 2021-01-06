package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

const helpText = "web3-cli - jsonrpc command line interface\n\nUsage: web3-cli method [param...]"

func main() {
	log.SetFlags(0)

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
					log.Printf("value %s should convert to number\n", p)
					return
				}
			case p == "true" || p == "false":
				if v, err := strconv.ParseBool(p); err == nil {
					params = append(params, v)
				} else {
					log.Printf("value %s should convert to bool\n", p)
					return
				}
			default:
				params = append(params, p)
			}
		}
	}

	endpoint := os.Getenv("web3")
	if endpoint == "" {
		fmt.Println("no env setted,use default localhost web3 server")
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

	log.Printf("Calling %s with %v", method, params)
	var result interface{}
	if err := ethclient.CallContext(ctx, &result, method, params...); err != nil {
		log.Printf("Call failed: %s\n", err)
		return
	}

	encoding := json.NewEncoder(os.Stdout)
	encoding.SetIndent("", "  ")
	_ = encoding.Encode(result)
}
