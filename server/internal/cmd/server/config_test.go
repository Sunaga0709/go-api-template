//go:build !integration

package main

import (
	"reflect"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    config
		wantErr bool
	}{
		// 正常系
		{
			name: "all env set",
			env: map[string]string{
				"SERVER_PORT":        "8080",
				"ENVIRONMENT":        "local",
				"LOG_LEVEL":          "debug",
				"HIDDEN_SUCCESS_LOG": "false",
			},
			want: config{
				Port:             "8080",
				Environment:      "local",
				LogLevel:         "debug",
				HiddenSuccessLog: false,
			},
			wantErr: false,
		},
		{
			name: "log level and hidden success log default",
			env: map[string]string{
				"SERVER_PORT": "8080",
				"ENVIRONMENT": "prd",
			},
			want: config{
				Port:             "8080",
				Environment:      "prd",
				LogLevel:         "info",
				HiddenSuccessLog: true,
			},
			wantErr: false,
		},
		// 異常系
		{
			name: "missing SERVER_PORT",
			env: map[string]string{
				"ENVIRONMENT": "local",
			},
			want:    config{},
			wantErr: true,
		},
		{
			name: "missing ENVIRONMENT",
			env: map[string]string{
				"SERVER_PORT": "8080",
			},
			want:    config{Port: "8080"},
			wantErr: true,
		},
		{
			name:    "missing all required",
			env:     map[string]string{},
			want:    config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnvForTest(t, tt.env)
			got, err := loadConfig()
			if (err != nil) != tt.wantErr {
				t.Fatalf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
