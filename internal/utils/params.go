package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

func parseJsonRawMsg(args []json.RawMessage) (params []any, err error) {
	for _, v := range args {
		switch {
		case bytes.HasPrefix(v, []byte("{")):
			var raw map[string]any
			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return nil, fmt.Errorf("invalid json object %s", v)
			}
			params = append(params, raw)
		case bytes.HasPrefix(v, []byte("[")):
			var raw []json.RawMessage
			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return nil, fmt.Errorf("invalid json array %s", v)
			}
			var got []any
			for _, item := range raw {
				res, err := parseJsonRawMsg([]json.RawMessage{item})
				if err != nil {
					return nil, err
				}
				got = append(got, res...)
			}
			params = append(params, got)
		case IsNumber(string(v)):
			if t, ok := new(big.Int).SetString(string(v), 10); ok {
				params = append(params, "0x"+t.Text(16))
			}
		default:
			var res any
			if err := json.Unmarshal(v, &res); err != nil {
				return nil, err
			}
			params = append(params, res)
		}
	}
	return params, nil
}

func ParseArgs(args []string) (params []any, err error) {
	for _, v := range args {
		switch {
		case IsNumber(v):
			if v, ok := new(big.Int).SetString(v, 10); ok {
				params = append(params, "0x"+v.Text(16))
				continue
			}
			fallthrough
		case v == "false", v == "true":
			v, _ := strconv.ParseBool(v)
			params = append(params, v)
		case v == "null":
			params = append(params, nil)
		case strings.HasPrefix(v, "["):
			var raw []json.RawMessage
			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return nil, fmt.Errorf("invalid json array %s", v)
			}
			var got []any
			for _, item := range raw {
				res, err := parseJsonRawMsg([]json.RawMessage{item})
				if err != nil {
					return nil, err
				}
				got = append(got, res...)
			}
			params = append(params, got)
		case strings.HasPrefix(v, "{"):
			var raw map[string]any
			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return nil, fmt.Errorf("invalid json object %s", v)
			}
			params = append(params, raw)
		default:
			params = append(params, v)
		}
	}
	return params, nil
}

func PrintJson(w io.Writer, data any, pretty bool) error {
	jsonEncoder := json.NewEncoder(w)
	if pretty {
		jsonEncoder.SetIndent("", "    ")
	}
	return jsonEncoder.Encode(data)
}
