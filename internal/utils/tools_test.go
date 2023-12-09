package utils

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestDecodeRawTransaction(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    string
		wantErr bool
	}{
		{
			name:    "0x4595df8f97c1b69b295a2ac9d2a1b95df6c876f7b249147e4f5ae3d34f812119",
			args:    "0x02f871018301f215808505ea45600082565f94388c818ca8b9251b393131c08a736a67ccb1929787dad10aed6a5f0d80c080a04bd8ec5c0a0bb85f2b1f02b62411947334912704374f1a51fb02332a87e17503a0291fa7c0f4295b0a018daf87ff8e7464db454da3147a2457d3576ff846ddceb4",
			want:    `{"type":"0x2","chainId":"0x1","nonce":"0x1f215","to":"0x388c818ca8b9251b393131c08a736a67ccb19297","gas":"0x565f","gasPrice":null,"maxPriorityFeePerGas":"0x0","maxFeePerGas":"0x5ea456000","value":"0xdad10aed6a5f0d","input":"0x","accessList":[],"v":"0x0","r":"0x4bd8ec5c0a0bb85f2b1f02b62411947334912704374f1a51fb02332a87e17503","s":"0x291fa7c0f4295b0a018daf87ff8e7464db454da3147a2457d3576ff846ddceb4","yParity":"0x0","hash":"0x4595df8f97c1b69b295a2ac9d2a1b95df6c876f7b249147e4f5ae3d34f812119"}`,
			wantErr: false,
		},
		{
			name:    "0x942c20100c89fbadf059b556a54fb89d3ef0229fd2aac82c49433af2cff14f22",
			args:    "f86b1c85028fa6ae008252089410dbe048cad44a462ba4b3c3d4f0975efd49206f87470de4df8200008025a07b9ba0158500031108d4fa7817314aeef01642d8c398c3118144dc762dc7cc58a05a720d7385e977ca12043b1b45d950ed2f1f0c5436a75fc4d22c1b2a24876dec",
			want:    `{"type":"0x0","chainId":"0x1","nonce":"0x1c","to":"0x10dbe048cad44a462ba4b3c3d4f0975efd49206f","gas":"0x5208","gasPrice":"0x28fa6ae00","maxPriorityFeePerGas":null,"maxFeePerGas":null,"value":"0x470de4df820000","input":"0x","v":"0x25","r":"0x7b9ba0158500031108d4fa7817314aeef01642d8c398c3118144dc762dc7cc58","s":"0x5a720d7385e977ca12043b1b45d950ed2f1f0c5436a75fc4d22c1b2a24876dec","hash":"0x942c20100c89fbadf059b556a54fb89d3ef0229fd2aac82c49433af2cff14f22"}`,
			wantErr: false,
		},
		{
			name:    "invalid hex string",
			args:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			wantErr: true,
		},
		{
			name:    "invalid transaction",
			args:    "01fd477a485e3c3072b5f84a5c3ab9237c121f269853b0b430ea212d76d7299d",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeRawTransaction(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("DecodeRawTransaction() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNewRandomAddress(t *testing.T) {
	type args struct {
		prefix string
		suffix string
		thread int
	}
	tests := []struct {
		name  string
		args  args
		args2 args
	}{
		{"1", args{"0x0", "0x00", 1}, args{"0", "00", 1}},
		{"2", args{"A", "BA", 4}, args{"a", "ba", 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomAddress(tt.args.prefix, tt.args.suffix, tt.args.thread)
			addr := hex.EncodeToString(got.Address.Bytes())
			if !strings.HasPrefix(addr, tt.args2.prefix) || !strings.HasSuffix(addr, tt.args2.suffix) {
				t.Errorf("TestNewRandomAddress() = %s", addr)
			}
		})
	}
}
