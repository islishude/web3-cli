package utils

import (
	"math/big"
	"strings"
)

func ToWei(amount float64) *big.Int {
	eth := big.NewFloat(amount)
	eth = eth.Mul(eth, big.NewFloat(1e18))
	wei, _ := eth.Int(new(big.Int))
	return wei
}

func ToBigInt(s string) *big.Int {
	// trim quotation mark
	s = strings.Trim(s, `"`)

	// no negative number
	if strings.HasPrefix(s, "-") {
		return nil
	}

	// no decimals
	if strings.ContainsRune(s, '.') {
		return nil
	}

	// remove underscores and convert to lower case
	s = strings.Join(strings.Split(strings.ToLower(s), "_"), "")
	if strings.HasPrefix(s, "0x") {
		s = strings.TrimPrefix(s, "0x")
		i, _ := new(big.Int).SetString(s, 16)
		return i
	}

	if strings.HasPrefix(s, "0b") {
		s = strings.TrimPrefix(s, "0b")
		i, _ := new(big.Int).SetString(s, 2)
		return i
	}

	if strings.HasPrefix(s, "0o") {
		s = strings.TrimPrefix(s, "0o")
		i, _ := new(big.Int).SetString(s, 8)
		return i
	}

	if strings.Contains(s, "e") {
		flt, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
		if err == nil {
			i, _ := flt.Int(new(big.Int))
			return i
		}
	}

	i, _ := new(big.Int).SetString(s, 10)
	return i
}
