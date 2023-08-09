package abis

import (
	"embed"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	//go:embed abi/*.json
	abis embed.FS
)

var abiSet = make(map[string]*abi.ABI)

func init() {
	const dir = "abi"
	entries, err := abis.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, item := range entries {
		fileName := item.Name()

		file, err := abis.Open(filepath.Join(dir, fileName))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		parsed, err := abi.JSON(file)
		if err != nil {
			panic(err)
		}
		name := fileName[:len(fileName)-len(filepath.Ext(fileName))]
		abiSet[strings.ToLower(name)] = &parsed
	}
}
