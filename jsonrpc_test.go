package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/islishude/web3-cli/internal/utils"
	"github.com/urfave/cli/v2"
)

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type jsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  []any           `json:"params,omitempty"`
	Result  any             `json:"result,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
}

func TestJsonrpcCall(t *testing.T) {
	tests := []struct {
		jsonrpcMethod string
		jsonrpcArgs   []string
		parsedArgs    []any
		jsonrpcResult any
		wantErr       bool
		err           *jsonError
	}{
		{jsonrpcMethod: "rpc_1", jsonrpcArgs: []string{"true", "1"}, parsedArgs: []any{true, "0x1"}, jsonrpcResult: "ok"},
		{jsonrpcMethod: "rpc_2", jsonrpcArgs: []string{`["hello", "world"]`}, parsedArgs: []any{[]any{"hello", "world"}}, jsonrpcResult: []float64{1, 2, 3}},
		{jsonrpcMethod: "rpc_3", jsonrpcArgs: []string{}, parsedArgs: []any{}, wantErr: true, err: &jsonError{Message: "error message"}},
		{jsonrpcMethod: "rpc_4", jsonrpcArgs: []string{"null", `{"key":"value"}`}, parsedArgs: []any{nil, map[string]any{"key": "value"}}, jsonrpcResult: map[string]any{"hash": "0xbeaf"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.jsonrpcMethod, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var msg jsonrpcMessage

				if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
					_ = json.NewEncoder(w).Encode(&jsonrpcMessage{
						Version: msg.Version,
						ID:      msg.ID,
						Error:   &jsonError{Message: "invalid json payload"},
					})
					return
				}

				var respMsg = &jsonrpcMessage{Version: msg.Version, ID: msg.ID}

				if msg.Method != tt.jsonrpcMethod {
					respMsg.Error = &jsonError{Message: fmt.Sprintf("jsonrpc method name is not the same: expect %s but got %s", tt.jsonrpcMethod, msg.Method)}
					t.Errorf("%+v", respMsg.Error.Message)
					_ = json.NewEncoder(w).Encode(respMsg)
					return
				}

				if len(msg.Params) != len(tt.jsonrpcArgs) {
					respMsg.Error = &jsonError{Message: "jsonrpc param length is not the same with expected"}
					t.Errorf("%+v", respMsg.Error.Message)
					_ = json.NewEncoder(w).Encode(respMsg)
					return
				}

				for i, v := range tt.parsedArgs {
					if p := msg.Params[i]; !reflect.DeepEqual(v, p) {
						respMsg.Error = &jsonError{Message: fmt.Sprintf("jsonrpc param index %d is not the same: typeGot %v typeWant %v",
							i, reflect.TypeOf(p).Kind(), reflect.TypeOf(v).Kind())}
						t.Errorf("%+v", respMsg.Error.Message)
						_ = json.NewEncoder(w).Encode(respMsg)
						return
					}
				}

				if tt.wantErr {
					respMsg.Error = tt.err
					_ = json.NewEncoder(w).Encode(respMsg)
					return
				}
				respMsg.Result = tt.jsonrpcResult
				_ = json.NewEncoder(w).Encode(respMsg)
			}))
			defer ts.Close()

			set := flag.NewFlagSet(tt.jsonrpcMethod, 0)
			_ = set.Parse(append([]string{tt.jsonrpcMethod}, tt.jsonrpcArgs...))

			cctx := cli.NewContext(nil, set, nil)

			rpc, err := rpc.DialContext(cctx.Context, ts.URL)
			if err != nil {
				t.Errorf("can not connect the test server: %s", err)
				return
			}

			logger := new(bytes.Buffer)
			callErr := JsonrpcCall(cctx, rpc, logger)
			if (callErr != nil) != tt.wantErr {
				t.Errorf("JsonrpcCall() error = %v, wantError %v", callErr, tt.wantErr)
			}

			if tt.wantErr {
				if tt.err == nil || callErr == nil || callErr.Error() != tt.err.Message {
					t.Errorf("JsonrpcCall() gotError = %v, wantError %v", callErr, tt.err.Message)
				}
			}

			if callErr == nil {
				wantBuffer := new(bytes.Buffer)
				_ = utils.PrintJson(wantBuffer, tt.jsonrpcResult, true)
				if want, got := wantBuffer.Bytes(), logger.Bytes(); !bytes.Equal(want, got) {
					t.Errorf("JsonrpcCall() output = %s, expected %s", got, want)
				}
			}
		})
	}
}
