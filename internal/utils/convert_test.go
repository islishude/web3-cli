package utils

import (
	"math/big"
	"reflect"
	"testing"
)

func TestToBigInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{`invalid`, args{`xyz`}, nil},
		{`1e4`, args{"1e4"}, big.NewInt(1e4)},
		{`1e18`, args{"1e18"}, big.NewInt(1e18)},
		{`1`, args{"1"}, big.NewInt(1)},
		{`-1`, args{`"-1"`}, nil},
		{`0x1`, args{`0x1`}, big.NewInt(1)},
		{`0b1`, args{`0b1`}, big.NewInt(1)},
		{`0o1`, args{`0o1`}, big.NewInt(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBigInt(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBigInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToWei(t *testing.T) {
	type args struct {
		amount float64
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{"0.1", args{0.1}, big.NewInt(1e17)},
		{"0.01", args{0.01}, big.NewInt(1e16)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToWei(tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToWei() = %v, want %v", got, tt.want)
			}
		})
	}
}
