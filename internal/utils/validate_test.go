package utils

import (
	"testing"
)

func TestIsAddress(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"valid", args{"0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"}, true},
		{"invalid-1", args{"0xabc"}, false},
		{"invalid-2", args{"0xzy"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAddress(tt.args.v); got != tt.want {
				t.Errorf("IsAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"0xdEa01Dbeaf23", true},
		{"123abcdef", true},
		{"0xqwert", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsHex(tt.name); got != tt.want {
				t.Errorf("IsHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
