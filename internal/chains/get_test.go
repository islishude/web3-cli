package chains

import (
	"reflect"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	chainNameSet := make(map[string]bool)
	for _, chain := range Buintin {
		if strings.ToLower(chain.Name) != chain.Name {
			t.Errorf("chain %s name should be lower format", chain.Name)
		}
		if chainNameSet[chain.Name] {
			t.Errorf("chain %s name is duplicated", chain.Name)
		}
		chainNameSet[chain.Name] = true

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

	tests := []struct {
		name string
		arg  string
		want *Chain
	}{
		{"unknown-0", "", nil},
		{"unknown-1", "unknown", nil},
		{"unknown-2", "3", nil},
		{"eth-0", "mainnet", ethMainnetChain},
		{"eth-1", "eth-mainnet", ethMainnetChain},
		{"eth-2", "eth", ethMainnetChain},
		{"eth-3", "1", ethMainnetChain},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
