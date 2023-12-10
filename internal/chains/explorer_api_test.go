package chains

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/islishude/web3-cli/internal/utils"
)

func TestExplorer(t *testing.T) {
	// api key is not required
	apiKeyList := map[string]string{
		"eth":     "ENV_ETHERSCAN_API_KEY",
		"goerli":  "ENV_ETHERSCAN_API_KEY",
		"sepolia": "ENV_ETHERSCAN_API_KEY",
		"holesky": "ENV_ETHERSCAN_API_KEY",
	}

	addrList := map[string]string{
		// https://etherscan.io/contractsVerified
		"eth": "0xdac17f958d2ee523a2206206994597c13d831ec7",
		// https://goerli.etherscan.io/contractsVerified
		"goerli": "0x0E641aeAB50481B521ce3051cDb38a2D9ac9C9cc",
		// https://sepolia.etherscan.io/contractsVerified
		"sepolia": "0xaaA68C69e625d349c734318423eB07bA1dC1101D",
		// https://holesky.etherscan.io/contractsVerified
		"holesky": "0xb7fb99e86f93dc3047a12932052236d853065173",
		// https://arbiscan.io/contractsVerified
		"arbitrum": "0x1a3b50Bd09594f96dDC192396CE41256EBe0726e",
		// https://goerli.arbiscan.io/contractsVerified
		"arbitrum-goerli": "0x1d6032b5e4044d1F870D65f56C067d3D79CbE35d",
		// https://optimistic.etherscan.io/contractsVerified
		"op": "0xb54b6262407DB03Be70B6e697DE8F880c94a1461",
		// https://goerli-optimism.etherscan.io/contractsVerified
		"op-goerli": "0x98108048528e45f621ad024a1a36f572fcbc64cb",
		// https://bscscan.com/contractsVerified
		"bsc": "0xFf179bea631ef161aACe1E0764e7587c0541a566",
		// https://testnet.bscscan.com/contractsVerified
		"bsc-testnet": "0x7bf307427fee18d2d075b175aa2e3d8653834e0d",
		// https://polygonscan.com/contractsVerified
		"polygon": "0x72271554CC05C22CAf18e2088f2A6d6e2b3a52e1",
		// https://mumbai.polygonscan.com/contractsVerified
		"polygon-mubai": "0x9eF761219423A06752F8552c6033AB7B07bd6122",
		// https://basescan.org/contractsVerified
		"base": "0x7485d64788a1e8609199ff1d06f7075c56f9f22d",
		// https://goerli.basescan.org/contractsVerified
		"base-goerli": "0x87f0E0922207C3F1b88df4a78b8e3bd83C7C7A1F",
		// https://explorer.metis.io/contractsverified
		"metis": "0xEA32A96608495e54156Ae48931A7c20f0dcc1a21",
		// https://goerli.explorer.metisdevops.link/verified-contracts
		"metis-goerli": "0xfAf6E5c9463334D154f39596764493576187a5d2",
	}

	var validIt = func(item *Chain) error {
		addr, hasAddr := addrList[item.Name]
		if !hasAddr {
			return fmt.Errorf("address undefined")
		}

		if !utils.IsAddress(addr) {
			return fmt.Errorf("invalid address")
		}

		apiKey := apiKeyList[item.Name]
		if strings.HasPrefix(apiKey, "ENV") {
			apiKey = os.Getenv(apiKey)
		}

		explorer, err := utils.URLToGetABI(item.Explorer, addr, apiKey)
		if err != nil {
			return fmt.Errorf("invalid explorer endpoint: %s", err)
		}

		// test if the address exists
		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			jsonrpc, err := ethclient.DialContext(ctx, item.Endpoint)
			if err != nil {
				return fmt.Errorf("failed to connect to rpc %s: %s", item.Endpoint, err)
			}
			defer jsonrpc.Close()

			code, err := jsonrpc.CodeAt(ctx, common.HexToAddress(addr), nil)
			if err != nil {
				return fmt.Errorf("failed to get code for %s: %s", addr, err)
			}

			if len(code) == 0 {
				return fmt.Errorf("address %s is not a contract", addr)
			}
		}

		// fetch abi
		{
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			expReq, err := http.NewRequestWithContext(ctx, http.MethodGet, explorer, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %s", err)
			}

			expResp, err := http.DefaultClient.Do(expReq)
			if err != nil {
				return fmt.Errorf("failed to call the explorer: %s", err)
			}
			defer expResp.Body.Close()

			var res map[string]string
			if err := json.NewDecoder(expResp.Body).Decode(&res); err != nil {
				return fmt.Errorf("failed to decode result to json: %s", err)
			}
			if res["status"] == "0" {
				return fmt.Errorf("error from explorer: %s: %s", res["message"], res["result"])
			}

			if _, err := abi.JSON(strings.NewReader(res["result"])); err != nil {
				return fmt.Errorf("invalid abi json: %s", err)
			}
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
				t.Log(err)
			}
		})
	}
}
