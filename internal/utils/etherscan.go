package utils

import (
	"fmt"
	"net/url"
)

func URLForExpABI(base, contAddr, key string) (string, error) {
	got, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	if schema := got.Scheme; schema != "http" && schema != "https" {
		return "", fmt.Errorf("not an http schema url")
	}

	queryString := got.Query()
	queryString.Add("address", contAddr)
	queryString.Add("module", "contract")
	queryString.Add("action", "getabi")
	if key != "" {
		queryString.Add("apikey", key)
	}
	got.RawQuery = queryString.Encode()
	return got.String(), nil
}
