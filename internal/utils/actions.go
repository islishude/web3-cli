package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/islishude/bigint"
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
		var v bigint.Int
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return fmt.Errorf("invalid decimal number: %s", s)
		}
		return ctx.Set(flagName, v.ToInt().String())
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
	case strings.HasSuffix(s, "gwei"):
		v, err := strconv.ParseFloat(strings.TrimLeft(s, "gwei"), 64)
		if err != nil {
			return fmt.Errorf("invalid gwei value: %s", s)
		}
		return ctx.Set(flagName, (*hexutil.Big)(ToGWei(v)).String())
	default:
		var v bigint.Int
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return fmt.Errorf("invalid decimal number: %s", s)
		}
		return ctx.Set(flagName, v.ToInt().String())
	}
}
