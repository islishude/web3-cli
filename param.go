package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

func parseArgs(args []string) ([]interface{}, error) {
	var params []interface{}
	for _, p := range args {
		switch {
		case regexp.MustCompile(`^[0-9]+$`).MatchString(p):
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
