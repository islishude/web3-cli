package chains

import (
	"strconv"
	"strings"

	"github.com/islishude/web3-cli/internal/utils"
)

type Chain struct {
	Name     string
	Id       uint64
	Endpoint string   `json:",omitempty"`
	Explorer string   `json:",omitempty"`
	Alias    []string `json:",omitempty"`
}

func Get(n string) *Chain {
	if n == "" {
		return nil
	}

	if utils.IsNumber(n) {
		if s, err := strconv.ParseUint(n, 10, 32); err == nil {
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
