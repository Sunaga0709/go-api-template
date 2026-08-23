//go:build !integration

package main

import (
	"reflect"
	"testing"
)

func Test_loadSecret(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    secret
		wantErr bool
	}{
		// 正常系
		{
			name: "all env set",
			env: map[string]string{
				"DATABASE_URL": "postgres://user:password@localhost:5432/app?sslmode=disable",
			},
			want: secret{
				DatabaseURL: "postgres://user:password@localhost:5432/app?sslmode=disable",
			},
			wantErr: false,
		},
		// 異常系
		{
			name:    "missing DATABASE_URL",
			env:     map[string]string{},
			want:    secret{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnvForTest(t, tt.env)
			got, err := loadSecret()
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadSecret() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}
