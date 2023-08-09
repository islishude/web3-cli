package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/urfave/cli/v2"
)

func EthCallFromFlagAction(ctx *cli.Context, s string) error {
	if !IsAddress(s) {
		return fmt.Errorf("invalid address: %s", s)
	}
	return nil
}

func EthCallToFlagAction(ctx *cli.Context, s string) error {
	return EthCallFromFlagAction(ctx, s)
}

func EthCallHeightFlagAction(ctx *cli.Context, s string) error {
	const flagName = "call-height"

	switch s {
	case "safe", "finalized", "latest", "earliest", "pending":
		return nil
	default:
		if v := ToBigInt(s); v != nil {
			return ctx.Set(flagName, v.String())
		}
		return fmt.Errorf("invalid number: %s", s)
	}
}

func EthCallValueFlagAction(ctx *cli.Context, s string) error {
	const flagName = "call-value"
	switch {
	case strings.HasSuffix(s, "eth"):
		v, err := strconv.ParseFloat(strings.TrimLeft(s, "eth"), 64)
		if err != nil {
			return fmt.Errorf("invalid ether value: %s", s)
		}
		return ctx.Set(flagName, (*hexutil.Big)(ToWei(v)).String())
	default:
		if v := ToBigInt(s); v != nil {
			return ctx.Set(flagName, v.String())
		}
		return fmt.Errorf("invalid number: %s", s)
	}
}
