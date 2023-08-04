package utils

import (
	"net/url"
	"reflect"
	"testing"
)

func TestURLForExpABI(t *testing.T) {
	type args struct {
		base     string
		contAddr string
		key      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{}, "", true},
		{"wrong url", args{" http://foo.html", "", ""}, "", true},
		{"not an http schema", args{"s3://test/api", "", ""}, "", true},
		{"with key", args{"https://test.com/api", "0x123", "thekey"}, "https://test.com/api?address=0x123&module=contract&action=getabi&apikey=thekey", false},
		{"without key", args{"https://test.com/api", "0xabc", ""}, "https://test.com/api?address=0xabc&module=contract&action=getabi", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := URLForExpABI(tt.args.base, tt.args.contAddr, tt.args.key)
			if hasErr := err != nil; hasErr {
				if hasErr != tt.wantErr {
					t.Errorf("URLForExpABI() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			got1, err := url.Parse(got)
			if err != nil {
				t.Fatal(err)
			}

			want1, err := url.Parse(tt.want)
			if err != nil {
				t.Fatal(err)
			}

			if reflect.DeepEqual(got1, want1) {
				t.Errorf("URLForExpABI() = %v, want %v", got, tt.want)
			}
		})
	}
}
