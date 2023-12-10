package abis

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Unpack(abiIns *abi.ABI, methodName string, output []byte) ([]any, error) {
	methodIns, ok := abiIns.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s does not exist", methodName)
	}
	val, err := abiIns.Unpack(methodIns.Name, output)
	if err != nil {
		return nil, err
	}

	res := make([]any, 0, len(val))
	for idx, item := range val {
		res = append(res, unpack(item, &methodIns.Outputs[idx].Type))
	}

	return res, nil
}

func unpack(item any, typ *abi.Type) any {
	refVal := reflect.ValueOf(item)

	switch refVal.Kind() {
	case reflect.Array:
		switch typ.T {
		case abi.ArrayTy:
			var res []any
			for i := 0; i < refVal.Len(); i++ {
				res = append(res, unpack(refVal.Index(i).Interface(), typ.Elem))
			}
			return res
		case abi.FixedBytesTy:
			res := make([]byte, 0, refVal.Len())
			for i := 0; i < refVal.Len(); i++ {
				res = append(res, refVal.Index(i).Interface().(byte))
			}
			return hexutil.Bytes(res).String()
		}
	case reflect.Slice:
		switch typ.T {
		case abi.BytesTy:
			res := make([]byte, 0, refVal.Len())
			for i := 0; i < refVal.Len(); i++ {
				res = append(res, refVal.Index(i).Interface().(byte))
			}
			return hexutil.Bytes(res).String()
		default:
			var res []any
			for i := 0; i < refVal.Len(); i++ {
				res = append(res, unpack(refVal.Index(i).Interface(), typ.Elem))
			}
			return res
		}
	case reflect.Struct:
		res := make(map[string]any)
		refType := refVal.Type()
		for i := 0; i < refType.NumField(); i++ {
			res[refType.Field(i).Name] = unpack(refVal.Field(i).Interface(), typ.TupleElems[i])
		}
		return res
	}

	switch v := item.(type) {
	case *big.Int:
		return (*hexutil.Big)(v).String()
	case common.Address:
		return v.Hex()
	case uint8, uint16, uint32, uint64:
		return (*hexutil.Big)(new(big.Int).SetUint64(refVal.Convert(reflect.TypeOf(uint64(0))).Interface().(uint64))).String()
	case int8, int16, int32, int64:
		return (*hexutil.Big)(big.NewInt(refVal.Convert(reflect.TypeOf(int64(0))).Interface().(int64))).String()
	default:
		return v
	}
}
