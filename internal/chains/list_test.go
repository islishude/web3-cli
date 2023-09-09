package chains

import (
	"strings"
	"testing"
)

func TestIfChainConfigDuplicated(t *testing.T) {
	chainNameSet, chainIdSet := make(map[string]bool), make(map[uint64]bool)
	rpcSet, explorerSet := make(map[string]bool), make(map[string]bool)

	for _, chain := range Buintin {
		if strings.ToLower(chain.Name) != chain.Name {
			t.Errorf("chain %s name should be lower case", chain.Name)
		}

		if chainNameSet[chain.Name] {
			t.Errorf("chain %s name is duplicated", chain.Name)
		}
		chainNameSet[chain.Name] = true

		if chainIdSet[chain.Id] {
			t.Errorf("chain id %s is duplicated", chain.Name)
		}
		chainIdSet[chain.Id] = true

		if rpcSet[chain.Endpoint] {
			t.Errorf("chain endpoint %s is duplicate", chain.Endpoint)
		}
		rpcSet[chain.Endpoint] = true

		if explorerSet[chain.Explorer] {
			t.Errorf("warning: chain explorer %s is duplicate", chain.Endpoint)
		}
		rpcSet[chain.Explorer] = true

		if len(chain.Alias) == 0 && chain.Alias != nil {
			t.Errorf("no alias for %s, change it to nil", chain.Name)
		}

		for _, item := range chain.Alias {
			if strings.ToLower(item) != item {
				t.Errorf("alias %s for chain %s name should be lower format", item, chain.Name)
			}

			if chainNameSet[item] {
				t.Errorf("alias %s for chain %s name is duplicated", item, chain.Name)
			}
			chainNameSet[item] = true
		}
	}
}
