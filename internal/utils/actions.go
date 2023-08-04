package utils

import (
	"fmt"
	"math/big"
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
		switch {
		case IsHex(s):
			v, err := hexutil.DecodeBig(s)
			if err != nil {
				return err
			}
			return ctx.Set(flagName, (*hexutil.Big)(v).String())
		case IsNumber(s):
			v, ok := new(big.Int).SetString(s, 10)
			if ok {
				return fmt.Errorf("invalid decimal number: %s", s)
			}
			return ctx.Set(flagName, (*hexutil.Big)(v).String())
		default:
			return fmt.Errorf("invalid height value: %s", s)
		}
	}
}

func EthCallValueFlagAction(ctx *cli.Context, s string) error {
	const flagName = "call-value"
	switch {
	case IsHex(s):
		v, err := hexutil.DecodeBig(s)
		if err != nil {
			return err
		}
		return ctx.Set(flagName, (*hexutil.Big)(v).String())
	case IsNumber(s):
		v, ok := new(big.Int).SetString(s, 10)
		if ok {
			return fmt.Errorf("invalid decimal number: %s", s)
		}
		return ctx.Set(flagName, (*hexutil.Big)(v).String())
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
		return fmt.Errorf("invalid value: %s", s)
	}
}
