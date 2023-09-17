package chains

import (
	"strconv"
	"strings"

	"github.com/islishude/web3-cli/internal/utils"
)

type Chain struct {
	Name     string   // chain name
	Id       uint64   // chain id
	Endpoint string   `json:",omitempty"` // the rpc endpoint, https and wss supported
	Explorer string   `json:",omitempty"` // explorer api full endpoint
	Alias    []string `json:",omitempty"` // name alias, e.g. the name of eth mainnet is eth, you can also add alias like eth-mainnet
}

func Get(n string) *Chain {
	if n == "" {
		return nil
	}

	if utils.IsNumber(n) {
		if s, err := strconv.ParseUint(n, 10, 64); err == nil {
			for _, item := range Buintin {
				if item.Id == s {
					return item
				}
			}
			return nil
		}
	}

	n = strings.ToLower(n)
	for _, item := range Buintin {
		if item.Name == n {
			return item
		}
		for _, alias := range item.Alias {
			if alias == n {
				return item
			}
		}
	}

	return nil
}
