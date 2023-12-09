package utils

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
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

type AddressInfo struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    common.Address
}

func (n AddressInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"PrivateKey": hexutil.Encode(crypto.FromECDSA(n.PrivateKey)),
		"PublicKey":  hexutil.Encode(crypto.FromECDSAPub(n.PublicKey)),
		"Address":    n.Address.Hex(),
	})
}

func getNewAddress() (addr AddressInfo, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return AddressInfo{}, err
	}

	addr.PrivateKey = privateKey
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return AddressInfo{}, errors.New("unknown public key type")
	}
	addr.PublicKey = publicKey

	addr.Address = crypto.PubkeyToAddress(*publicKey)
	return
}

func NewRandomAddress(prefix, suffix string, thread int) AddressInfo {
	prefix = strings.ToLower(strings.TrimPrefix(prefix, "0x"))
	suffix = strings.ToLower(strings.TrimPrefix(suffix, "0x"))

	ch := make(chan AddressInfo, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < thread; i++ {
		go func() {
			var pt = true
			var st = true

			for {
				select {
				case <-ctx.Done():
					return
				default:
					got, err := getNewAddress()
					if err != nil {
						continue
					}

					addr := hex.EncodeToString(got.Address[:])
					if prefix != "" {
						pt = strings.HasPrefix(addr, prefix)
					}

					if suffix != "" {
						st = strings.HasSuffix(addr, suffix)
					}
					if pt && st {
						ch <- got
						return
					}
				}
			}
		}()
	}
	return <-ch
}
