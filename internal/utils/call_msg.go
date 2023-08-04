package utils

import "github.com/ethereum/go-ethereum/common/hexutil"

type CallMsg struct {
	From  string        `json:"from"`
	To    string        `json:"to"`
	Data  hexutil.Bytes `json:"data,omitempty"`
	Value *hexutil.Big  `json:"value,omitempty"`
	Gas   hexutil.Uint  `json:"gas,omitempty"`
}
