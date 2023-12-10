package tools

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
)

func DecodeRawTransaction(r string) ([]byte, error) {
	bytx, err := hex.DecodeString(strings.TrimPrefix(r, "0x"))
	if err != nil {
		return nil, err
	}

	var tx types.Transaction
	if err := tx.UnmarshalBinary(bytx); err != nil {
		return nil, err
	}

	return tx.MarshalJSON()
}
