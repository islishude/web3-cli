package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []interface{}
		wantErr bool
	}{
		{"null", []string{"null"}, []interface{}{nil}, false},
		{"invalid json", []string{`{"name":invalid}`}, nil, true},
		{"eth_getBlockByNumber", []string{"100", "true"}, []interface{}{"0x64", true}, false},
		{"eth_getBlockByNumber_2", []string{"latest", "false"}, []interface{}{"latest", false}, false},
		{"eth_eastimateGas", []string{`{"from": "1","to": "2"}`, "latest"}, []interface{}{map[string]interface{}{"from": "1", "to": "2"}, "latest"}, false},
		{"debug_traceTransaction", []string{"0x82c6040b89e79d136af7368f993c8fa5856d690be8bba5533ff807218f0d7292", `{"tracer": "callTracer"}`}, []interface{}{"0x82c6040b89e79d136af7368f993c8fa5856d690be8bba5533ff807218f0d7292", map[string]interface{}{"tracer": "callTracer"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintJson(t *testing.T) {
	type args struct {
		data   any
		pretty bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			"general",
			args{map[string]int{"a": 1}, false},
			fmt.Sprintln(`{"a":1}`),
			false,
		},
		{
			"pretty",
			args{map[string]int{"a": 1}, true},
			fmt.Sprintf(`{%s    "a": 1%s}%s`, "\n", "\n", "\n"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := PrintJson(w, tt.args.data, tt.args.pretty); (err != nil) != tt.wantErr {
				t.Errorf("PrintJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PrintJson() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
