package chains

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestChainId(t *testing.T) {
	var validIt = func(item *Chain) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		parsed, err := url.Parse(item.Endpoint)
		if err != nil {
			return fmt.Errorf("invalid rpc endpint: %s", item.Endpoint)
		}

		switch parsed.Scheme {
		case "https", "wss":
		case "http", "ws":
			return fmt.Errorf("rpc endpoint %s does not use TLS", item.Endpoint)
		default:
			// unix socket is disabled
			return fmt.Errorf("invalid rpc endpoint schema: %s", parsed.Scheme)
		}

		client, err := ethclient.DialContext(ctx, item.Endpoint)
		if err != nil {
			return fmt.Errorf("failed to connect to rpc %s: %s", item.Endpoint, err)
		}
		defer client.Close()

		chainId, err := client.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("failed to call eth_chainId: %s", err)
		}

		if !chainId.IsUint64() {
			return errors.New("got a non-uin64 chain id")
		}

		if v := chainId.Uint64(); item.Id != v {
			return fmt.Errorf("expected chainId %d, but got %d", item.Id, v)
		}

		return nil
	}

	for _, item := range Buintin {
		if item.Name == "local" {
			continue
		}

		item := item
		t.Run(item.Name, func(t *testing.T) {
			t.Parallel()
			if err := validIt(item); err != nil {
				t.Logf("invalid chain config: %s", err)
			}
		})
	}
}
