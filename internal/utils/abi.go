package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func ABIMethodPack(abiIns *abi.ABI, abiArgs []string) (hexutil.Bytes, error) {
	if len(abiArgs) == 0 {
		return nil, fmt.Errorf("no method name provide")
	}

	methodIns, ok := abiIns.Methods[abiArgs[0]]
	if !ok {
		return nil, fmt.Errorf("method %s does not exist", abiArgs[0])
	}

	if methodIns.Type != abi.Function {
		return nil, fmt.Errorf("%s is not a function", methodIns.Name)
	}

	abiArgs = abiArgs[1:]
	if len(methodIns.Inputs) != len(abiArgs) {
		return nil, errors.New("incorrect args number gived")
	}

	var abiParams []any
	for i := 0; i < len(methodIns.Inputs); i++ {
		abiArg, arg := methodIns.Inputs[i], abiArgs[i]

		switch abiArg.Type.T {
		case abi.IntTy, abi.UintTy:
			switch {
			case IsNumber(arg):
				val, ok := new(big.Int).SetString(arg, 10)
				if !ok {
					return nil, fmt.Errorf("param %s should be a valid number type", arg)
				}
				abiParams = append(abiParams, val)
			case IsHex(arg):
				val, err := hexutil.DecodeBig(arg)
				if err != nil {
					return nil, fmt.Errorf("param %s should be a valid hex number type(e.g. 0xabc)", arg)
				}
				abiParams = append(abiParams, val)
			default:
				return nil, fmt.Errorf("param %s should be a valid number type", arg)
			}
		case abi.BoolTy:
			val, err := strconv.ParseBool(arg)
			if err != nil {
				return nil, fmt.Errorf("param %s should be a valid bool type(e.g. true,false,1,0)", arg)
			}
			abiParams = append(abiParams, val)
		case abi.AddressTy:
			if !IsAddress(arg) {
				return nil, fmt.Errorf("param %s should be a valid address", arg)
			}
			abiParams = append(abiParams, common.HexToAddress(arg))
		case abi.StringTy:
			abiParams = append(abiParams, arg)
		case abi.BytesTy:
			val, err := hexutil.Decode(arg)
			if err != nil {
				return nil, fmt.Errorf("param %s should be a valid hex bytes(e.g. 0x123abc)", arg)
			}
			abiParams = append(abiParams, val)
		default:
			return nil, fmt.Errorf("not supported type %s", abiArg.Type.String())
		}
	}

	return abiIns.Pack(methodIns.Name, abiParams...)
}

func ABIMethodUnpack(abiIns *abi.ABI, methodName string, output []byte) ([]interface{}, error) {
	methodIns, ok := abiIns.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s does not exist", methodName)
	}
	return abiIns.Unpack(methodIns.Name, output)
}
