package abis

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func Unpack(abiIns *abi.ABI, methodName string, output []byte) ([]interface{}, error) {
	methodIns, ok := abiIns.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s does not exist", methodName)
	}
	return abiIns.Unpack(methodIns.Name, output)
}
