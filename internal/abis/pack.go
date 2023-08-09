package abis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/islishude/bigint"
	"github.com/islishude/web3-cli/internal/utils"
)

func Pack(abiIns *abi.ABI, abiArgs []string) ([]byte, error) {
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

	abiParams, err := pack(methodIns.Inputs, abiArgs)
	if err != nil {
		return nil, err
	}

	return abiIns.Pack(methodIns.Name, abiParams...)
}

func pack(argsIns abi.Arguments, args []string) (res []any, err error) {
	for i := 0; i < len(argsIns); i++ {
		abiArg, rawArg := argsIns[i], args[i]

		var enc any
		switch abiArg.Type.T {
		case abi.ArrayTy, abi.SliceTy:
			enc, err = sliceEncode(&abiArg.Type, rawArg)
		case abi.TupleTy:
			enc, err = tupleEncode(&abiArg.Type, rawArg)
		default:
			enc, err = nonDynamicEncode(&abiArg.Type, rawArg)
		}
		if err != nil {
			return nil, err
		}
		res = append(res, enc)
	}
	return res, nil
}

func sliceEncode(abiType *abi.Type, arg any) (any, error) {
	var rawArg []byte
	switch v := arg.(type) {
	case string:
		rawArg = []byte(v)
	case json.RawMessage:
		rawArg = v
	}

	var list []json.RawMessage
	err := json.Unmarshal(rawArg, &list)
	if err != nil {
		return nil, fmt.Errorf("invalid json input: %s", rawArg)
	}

	if abiType.T == abi.ArrayTy && len(list) != abiType.Size {
		return nil, fmt.Errorf("invalid length of the array type")
	}

	res := reflect.MakeSlice(reflect.SliceOf(abiType.Elem.GetType()), 0, len(list))
	for idx := 0; idx < len(list); idx++ {
		var val any
		switch abiType.Elem.T {
		case abi.ArrayTy, abi.SliceTy:
			val, err = sliceEncode(abiType.Elem, list[idx])
		case abi.TupleTy:
			val, err = tupleEncode(abiType.Elem, list[idx])
		default:
			val, err = nonDynamicEncode(abiType.Elem, list[idx])
		}
		if err != nil {
			return nil, err
		}
		res = reflect.Append(res, reflect.ValueOf(val))
	}
	return res.Interface(), nil
}

func tupleEncode(abiType *abi.Type, arg any) (res any, err error) {
	var rawArg []byte
	switch v := arg.(type) {
	case string:
		rawArg = []byte(v)
	case json.RawMessage:
		rawArg = v
	}

	var list []json.RawMessage
	if isBatch(rawArg) {
		if err := json.Unmarshal(rawArg, &list); err != nil {
			return nil, fmt.Errorf("invalid json input: %s", rawArg)
		}
		if len(list) != len(abiType.TupleElems) {
			return nil, fmt.Errorf("invalid length of the tuple type for %s", abiType.String())
		}
	} else {
		var tuple map[string]json.RawMessage
		if err := json.Unmarshal(rawArg, &tuple); err != nil {
			return nil, fmt.Errorf("invalid json input: %s", rawArg)
		}

		if len(tuple) != len(abiType.TupleElems) {
			return nil, fmt.Errorf("invalid length of the tuple type for %s", abiType.String())
		}

		for _, cname := range abiType.TupleRawNames {
			var item json.RawMessage
			item, has := tuple[abi.ToCamelCase(cname)]
			if !has {
				item, has = tuple[cname]
				if !has {
					return nil, fmt.Errorf("miss %s filed for tuple %s", cname, abiType.String())
				}
			}
			list = append(list, item)
		}
	}

	value := reflect.New(abiType.TupleType)
	for idx, item := range abiType.TupleElems {
		var val any
		switch item.T {
		case abi.ArrayTy, abi.SliceTy:
			val, err = sliceEncode(item, list[idx])
		case abi.TupleTy:
			val, err = tupleEncode(item, list[idx])
		default:
			val, err = nonDynamicEncode(item, list[idx])
		}
		if err != nil {
			return nil, err
		}
		if f := value.Elem().Field(idx); f.CanSet() {
			f.Set(reflect.ValueOf(val))
		} else {
			return nil, fmt.Errorf("tuple can't be set")
		}
	}
	return value.Elem().Interface(), nil
}

func nonDynamicEncode(abiArg *abi.Type, arg any) (any, error) {
	var rawArg string
	switch v := arg.(type) {
	case string:
		rawArg = v
	case json.RawMessage:
		if abiArg.T != abi.IntTy && abiArg.T != abi.UintTy {
			v = bytes.TrimLeft(v, `"`)
			v = bytes.TrimRight(v, `"`)
		}
		rawArg = string(v)
	}

	switch abiArg.T {
	case abi.IntTy, abi.UintTy:
		var i bigint.Int
		if err := json.Unmarshal([]byte(rawArg), &i); err != nil {
			return nil, fmt.Errorf("param %s should be a valid number type: %s", arg, err)
		}
		return i.ToInt(), nil
	case abi.BoolTy:
		val, err := strconv.ParseBool(rawArg)
		if err != nil {
			return nil, fmt.Errorf("param %s should be a valid bool type(e.g. true,false,1,0)", arg)
		}
		return val, nil
	case abi.StringTy:
		return arg, nil
	case abi.AddressTy:
		if !utils.IsAddress(rawArg) {
			return nil, fmt.Errorf("param %s is not a valid address", arg)
		}
		return common.HexToAddress(rawArg), nil
	case abi.BytesTy, abi.FixedBytesTy:
		val, err := hexutil.Decode(rawArg)
		if err != nil {
			return nil, fmt.Errorf("param %s should be a valid hex bytes(e.g. 0x123abc)", arg)
		}
		return val, nil
	default:
		return nil, fmt.Errorf("not supported type %s", abiArg.String())
	}
}
