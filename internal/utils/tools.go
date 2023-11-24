package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

type NewAddress struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func NewRandomAddress() (addr NewAddress, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return NewAddress{}, err
	}

	addr.PrivateKey = hexutil.Encode(crypto.FromECDSA(privateKey))

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return NewAddress{}, errors.New("unknown public key type")
	}
	addr.PublicKey = hexutil.Encode(crypto.FromECDSAPub(publicKey))

	addr.Address = crypto.PubkeyToAddress(*publicKey).Hex()
	return
}
