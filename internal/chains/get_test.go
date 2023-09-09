package chains

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want *Chain
	}{
		{"unknown-0", "", nil},
		{"unknown-1", "unknown", nil},
		{"unknown-2", "3", nil},
		{"eth-0", "mainnet", eth},
		{"eth-1", "eth-mainnet", eth},
		{"eth-2", "eth", eth},
		{"eth-3", "1", eth},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
