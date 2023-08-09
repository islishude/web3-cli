package abis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func Get(name string) *abi.ABI {
	if name == "" {
		return nil
	}

	return abiSet[strings.ToLower(name)]
}

func Builtin() (res []string) {
	for item := range abiSet {
		res = append(res, strings.ToUpper(item))
	}
	return
}

func Fetch(path string, isExplorerApi bool) (*abi.ABI, error) {
	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	switch pathURL.Scheme {
	case "http", "https":
		httpResp, err := http.Get(pathURL.String())
		if err != nil {
			return nil, err
		}
		defer httpResp.Body.Close()

		if isExplorerApi {
			// https://docs.etherscan.io/api-endpoints/contracts#get-contract-abi-for-verified-contract-source-codes
			var res map[string]string
			if err := json.NewDecoder(httpResp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("can not decode result from %s", pathURL.Host)
			}
			if res["status"] == "0" {
				return nil, fmt.Errorf("error from explorer: %s", res["message"])
			}
			reader = strings.NewReader(res["result"])
		} else {
			if httpResp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(httpResp.Body)
				return nil, fmt.Errorf("can not fetch abi from %s: %s", pathURL.Host, body)
			}
			reader = httpResp.Body
		}
	case "":
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	default:
		return nil, fmt.Errorf("%s is not supported to get abi", pathURL.Scheme)
	}

	parsed, err := abi.JSON(reader)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
