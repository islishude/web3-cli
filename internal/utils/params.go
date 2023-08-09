package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

func ParseArgs(args []string) (params []interface{}, err error) {
	for _, p := range args {
		switch {
		case IsNumber(p):
			if v, ok := new(big.Int).SetString(p, 10); ok {
				params = append(params, "0x"+v.Text(16))
				continue
			}
			fallthrough
		case p == "false", p == "true":
			v, _ := strconv.ParseBool(p)
			params = append(params, v)
		case p == "null":
			params = append(params, nil)
		case strings.HasPrefix(p, "["), strings.HasPrefix(p, "{"):
			var raw interface{}
			if err := json.Unmarshal([]byte(p), &raw); err != nil {
				return nil, fmt.Errorf("format of param %s is not correct: %v", p, err)
			}
			params = append(params, raw)
		default:
			params = append(params, p)
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
