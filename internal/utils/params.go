package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

func ParseArgs(args []string, useEthFormat bool) (params []any, err error) {
	for _, v := range args {
		switch {
		case IsNumber(v):
			if v, ok := new(big.Int).SetString(v, 10); ok {
				if useEthFormat {
					params = append(params, "0x"+v.Text(16))
				} else {
					if !v.IsInt64() {
						return nil, fmt.Errorf("number %s is too large", v.String())
					}
					params = append(params, v.Int64())
				}
				continue
			}
			fallthrough
		case v == "false", v == "true":
			v, _ := strconv.ParseBool(v)
			params = append(params, v)
		case v == "null":
			params = append(params, nil)
		case strings.HasPrefix(v, "["):
			var raw []any
			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return nil, fmt.Errorf("invalid json array %s", v)
			}
			params = append(params, raw)
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
