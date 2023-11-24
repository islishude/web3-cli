package utils

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestEthCallValueFlagAction(t *testing.T) {
	const flagName = "call-value"

	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			"eth-0", "0.001eth", "0x38d7ea4c68000", false,
		},
		{
			"eth-1", "0.1234 eth", "0x1b667a56d488000", false,
		},
		{
			"val-0", "16", "0x10", false,
		},
		{
			"val-1", "0x10", "0x10", false,
		},
		{
			"val-invalid-0", "test", "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet(tt.name, 0)
			set.String(flagName, tt.args, flagName)

			ctx := cli.NewContext(nil, set, nil)

			err := EthCallValueFlagAction(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EthCallValueFlagAction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if got := ctx.String(flagName); got != tt.want {
					t.Errorf("EthCallValueFlagAction() = %s, want %s", got, tt.want)
				}
			}
		})
	}
}

func TestEthCallHeightFlagAction(t *testing.T) {
	const flagName = "call-height"

	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{"latest", "latest", "latest", false},
		{"pending", "pending", "pending", false},
		{"0x10", "0x10", "0x10", false},
		{"16", "16", "0x10", false},
		{"invalid", "test", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet(tt.name, 0)
			set.String(flagName, tt.args, flagName)

			ctx := cli.NewContext(nil, set, nil)

			err := EthCallHeightFlagAction(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EthCallHeightFlagAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				if got := ctx.String(flagName); got != tt.want {
					t.Errorf("EthCallHeightFlagAction() = %s, want %s", got, tt.want)
				}
			}
		})
	}
}

func TestEthCallToFlagAction(t *testing.T) {
	const flagName = "test"

	tests := []struct {
		name    string
		args    string
		wantErr bool
	}{
		{"addr-0", "0xc72DD9CA4659fcBA254bb3858aCb8e572ec19aFe", false},
		{"addr-1", "0x620F6B61e9D5278172f06594f8Cc7FbFa89daF81", false},
		{"addr-2", "0x0000000000000000000000000000000000000000", false},
		{"addr-3", "144ce3c3146002ccd1de2124", true},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet(tt.name, 0)
			set.String(flagName, tt.args, flagName)

			ctx := cli.NewContext(nil, set, nil)

			if err := EthCallToFlagAction(ctx, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("EthCallToFlagAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEthCallGasFlagAction(t *testing.T) {
	const flagName = "call-gas"

	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			"val-0", "1", "1", false,
		},
		{
			"val-1", "0x10", "16", false,
		},
		{
			"invalid", "0.1", "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet(tt.name, 0)
			set.String(flagName, tt.args, flagName)

			ctx := cli.NewContext(nil, set, nil)

			err := EthCallGasFlagAction(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestEthCallGasFlagAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				if got := ctx.String(flagName); got != tt.want {
					t.Errorf("TestEthCallGasFlagAction() = %s, want %s", got, tt.want)
				}
			}
		})
	}
}
