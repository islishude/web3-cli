package tools

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestNewRandomAddress(t *testing.T) {
	type args struct {
		prefix string
		suffix string
		thread int
	}
	tests := []struct {
		name  string
		args  args
		args2 args
	}{
		{"1", args{"0x0", "0x00", 1}, args{"0", "00", 1}},
		{"2", args{"A", "BA", 4}, args{"a", "ba", 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomAddress(tt.args.prefix, tt.args.suffix, tt.args.thread)
			addr := hex.EncodeToString(got.Address.Bytes())
			if !strings.HasPrefix(addr, tt.args2.prefix) || !strings.HasSuffix(addr, tt.args2.suffix) {
				t.Errorf("TestNewRandomAddress() = %s", addr)
			}
		})
	}
}
