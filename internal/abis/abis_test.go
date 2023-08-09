package abis_test

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/islishude/web3-cli/internal/abis"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		isNil bool
	}{
		{"erc20", false},
		{"ERC20", false},
		{"erc721", false},
		{"", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := abis.Get(tt.name); tt.isNil && got != nil {
				t.Errorf("nil")
			}
		})
	}
}

//go:embed testdata/usdt.abi.json
var usdtAbi string

func TestFetch(t *testing.T) {
	parsed, err := abi.JSON(strings.NewReader(usdtAbi))
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/good" {
			fmt.Fprintln(w, usdtAbi)
			return
		}
		if r.URL.Path == "/bad" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "bad call")
			return
		}
		fmt.Fprintln(w, "not a valid json content")
	}))
	defer ts.Close()

	ests := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/good" {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "1", "result": usdtAbi})
			return
		}
		if r.URL.Path == "/bad" {
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "0", "message": "bad call"})
			return
		}
		fmt.Fprintln(w, "not a valid json content")
	}))
	defer ests.Close()

	type args struct {
		path          string
		isExplorerApi bool
	}
	tests := []struct {
		name    string
		args    args
		want    *abi.ABI
		wantErr bool
	}{
		{"local", args{filepath.Join(wd, "testdata/usdt.abi.json"), false}, &parsed, false},
		{"local, file not found", args{filepath.Join(wd, "/tmp/not-found.json"), false}, nil, true},
		{"unkown schema", args{"s3://abis/test.json", false}, nil, true},
		{"unknown remote host", args{fmt.Sprintf("https://unit-test-unknown-host-%d.dev", rand.Int()), false}, nil, true},
		{"good remote", args{ts.URL + "/good", false}, &parsed, false},
		{"bad remote", args{ts.URL + "/bad", false}, nil, true},
		{"bad json from remote", args{ts.URL, false}, nil, true},
		{"good etherscan", args{ests.URL + "/good", true}, &parsed, false},
		{"bad etherscan", args{ests.URL + "/bad", true}, nil, true},
		{"bad json from etherscan", args{ests.URL, true}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := abis.Fetch(tt.args.path, tt.args.isExplorerApi)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}
