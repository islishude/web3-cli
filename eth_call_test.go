package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/islishude/web3-cli/internal/chains"
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"

	_ "embed"
)

//go:embed internal/abis/abi/erc20.json
var erc20AbiRaw string

func TestContractCall(t *testing.T) {
	type jsonrpcMessage struct {
		Version string            `json:"jsonrpc,omitempty"`
		ID      json.RawMessage   `json:"id,omitempty"`
		Method  string            `json:"method,omitempty"`
		Params  []json.RawMessage `json:"params,omitempty"`
		Result  any               `json:"result,omitempty"`
		Error   *jsonError        `json:"error,omitempty"`
	}

	var callMegEqual = func(o1, o2 EthCallMsg) bool {
		if o1.From != o2.From {
			return false
		}

		if o1.To != o2.To {
			return false
		}

		if o1.Value == nil && o2.Value != nil {
			return false
		}

		if o1.Value.ToInt().Cmp(o2.Value.ToInt()) != 0 {
			return false
		}

		if !bytes.Equal(o2.Data, o1.Data) {
			return false
		}

		return o1.Gas == o2.Gas
	}

	tests := []struct {
		abiJson    string
		abiName    string
		abiMethod  string
		callParams EthCallMsg
		height     string
		abiArgs    []string
		wantInput  string
		result     []any
		abiOutput  hexutil.Bytes
		wantErr    bool
	}{
		{
			abiName:    "ERC20",
			abiMethod:  "balanceOf",
			callParams: EthCallMsg{To: "0xe339d1fc124f44ba910c6e171b79f48212011178"},
			height:     "latest",
			abiArgs:    []string{"0xb228200dfc46e5da38e751d586045f112b0eda2c"},
			wantInput:  "0x70a08231000000000000000000000000b228200dfc46e5da38e751d586045f112b0eda2c",
			result:     []any{big.NewInt(100)},
			abiOutput:  hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000064"),
		},
		{
			abiName:    "fetch-abi",
			abiJson:    erc20AbiRaw,
			abiMethod:  "totalSupply",
			callParams: EthCallMsg{To: "0x665e70200dc4c05cedbdd513e2218171e1ed0ad9", From: "0x7749df9ac4208558935515098c5050457124164e", Gas: hexutil.Uint(100), Value: (*hexutil.Big)(hexutil.MustDecodeBig("0x1"))},
			height:     "0x64",
			abiArgs:    []string{},
			wantInput:  "0x18160ddd",
			result:     []any{big.NewInt(100)},
			abiOutput:  hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000064"),
		},
		{
			abiName:    "ERC20",
			abiMethod:  "allowance",
			callParams: EthCallMsg{To: "0x922e95e7d1080e01390f246286b187981ba0ef0e"},
			height:     "latest",
			abiArgs:    []string{"0xf710d46ded2d2e056bf240c639599993b83e6d07", "0xb4aa3cfba2a1bbae9fe8ebd97e1a17fd78fb5416"},
			wantInput:  "0xdd62ed3e000000000000000000000000f710d46ded2d2e056bf240c639599993b83e6d07000000000000000000000000b4aa3cfba2a1bbae9fe8ebd97e1a17fd78fb5416",
			abiOutput:  hexutil.MustDecode("0x"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.abiMethod, func(t *testing.T) {
			tt := tt
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				if r.URL.Path == "/api" && tt.abiJson != "" {
					_ = json.NewEncoder(w).Encode(map[string]string{"status": "1", "result": tt.abiJson})
					return
				}

				var msg jsonrpcMessage
				if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
					msg.Error = &jsonError{Message: "invalid json payload for eth_call"}
					t.Errorf("%+v", msg.Error.Message)
					return
				}

				var respMsg = &jsonrpcMessage{Version: msg.Version, ID: msg.ID}

				// check method and params
				{
					if len(msg.Params) != 2 {
						respMsg.Error = &jsonError{Message: "the params should have 2 items"}
						t.Errorf("%+v", respMsg.Error.Message)
						_ = json.NewEncoder(w).Encode(respMsg)
						return
					}

					switch msg.Method {
					// check call msg
					case "eth_call":
						var callMsg EthCallMsg
						if err := json.Unmarshal(msg.Params[0], &callMsg); err != nil {
							respMsg.Error = &jsonError{Message: "the first parameter should be type CallMsg"}
							t.Errorf("%+v", respMsg.Error.Message)
							_ = json.NewEncoder(w).Encode(respMsg)
							return
						}

						tt.callParams.Data = hexutil.MustDecode(tt.wantInput)
						if !callMegEqual(callMsg, tt.callParams) {
							respMsg.Error = &jsonError{Message: fmt.Sprintf("the call msg is not the same: want %+v got %+v", tt.callParams, callMsg)}
							t.Errorf("%+v", respMsg.Error.Message)
							_ = json.NewEncoder(w).Encode(respMsg)
							return
						}

						var height string
						if err := json.Unmarshal(msg.Params[1], &height); err != nil {
							respMsg.Error = &jsonError{Message: "the second parameter should be a hex number for height"}
							t.Errorf("%+v", respMsg.Error.Message)
							_ = json.NewEncoder(w).Encode(respMsg)
							return
						}

						switch height {
						case "safe", "finalized", "latest", "earliest", "pending":
						default:
							if _, err := hexutil.DecodeBig(height); err != nil {
								respMsg.Error = &jsonError{Message: "the height parameter should be a valid hex string"}
								t.Errorf("%+v", respMsg.Error.Message)
								_ = json.NewEncoder(w).Encode(respMsg)
								return
							}
						}

						if height != tt.height {
							respMsg.Error = &jsonError{Message: "the height parameter is not the same"}
							t.Errorf("%+v", respMsg.Error.Message)
							_ = json.NewEncoder(w).Encode(respMsg)
							return
						}

					case "eth_getCode":
						var addr string
						if err := json.Unmarshal(msg.Params[0], &addr); err != nil {
							msg.Error = &jsonError{Message: "invalid address format for eth_getCode"}
							t.Errorf("%+v", msg.Error.Message)
							return
						}
						if addr != tt.callParams.To {
							msg.Error = &jsonError{Message: "invalid address"}
							t.Errorf("%+v", msg.Error.Message)
							return
						}
						respMsg.Result = tt.abiOutput
					}

					respMsg.Result = tt.abiOutput
					_ = json.NewEncoder(w).Encode(respMsg)
				}
			}))
			defer ts.Close()

			flagset := flag.NewFlagSet(tt.abiMethod, 0)
			flagset.String(ABINameFlag.Name, tt.abiName, "-abi-name")
			flagset.String(EthCallToFlag.Name, tt.callParams.To, "-call-to")

			if v := tt.callParams.From; v != "" {
				flagset.String(EthCallFromFlag.Name, v, "-call-from")
			}

			if v := tt.callParams.Value; v != nil {
				flagset.String(EthCallValueFlag.Name, v.String(), "-call-value")
			}

			if v := tt.callParams.Gas; v != 0 {
				flagset.String(EthCallGasFlag.Name, v.String(), "-call-gas")
			}

			if v := tt.height; v != "" {
				flagset.String(EthCallHeightFlag.Name, v, "-call-height")
			}

			if err := flagset.Parse(append([]string{tt.abiMethod}, tt.abiArgs...)); err != nil {
				t.Fatal(err)
				return
			}

			cctx := cli.NewContext(nil, flagset, nil)

			rpc, err := rpc.DialContext(cctx.Context, ts.URL)
			if err != nil {
				t.Errorf("can not connect the test server: %s", err)
				return
			}

			logger := new(bytes.Buffer)

			callErr := ContractCall(cctx, rpc, &chains.Chain{Name: "test", Endpoint: fmt.Sprintf("%s/rpc", ts.URL), Explorer: fmt.Sprintf("%s/api", ts.URL)}, logger)
			if (callErr != nil) != tt.wantErr {
				t.Errorf("ContractCall() error = %v, wantErr %v", callErr, tt.wantErr)
			}

			if callErr == nil {
				wantBuffer := new(bytes.Buffer)
				_ = utils.PrintJson(wantBuffer, tt.result, true)
				if want, got := wantBuffer.Bytes(), logger.Bytes(); !bytes.Equal(want, got) {
					t.Errorf("ContractCall() output = %s, expected %s", got, want)
				}
			}
		})
	}
}
