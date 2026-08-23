//go:build !integration

package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{level: ""},
			wantErr: false,
		},
		{
			name:    "trace",
			args:    args{level: "trace"},
			wantErr: false,
		},
		{
			name:    "Trace",
			args:    args{level: "Trace"},
			wantErr: false,
		},
		{
			name:    "TRACE",
			args:    args{level: "TRACE"},
			wantErr: false,
		},
		{
			name:    "tRACE",
			args:    args{level: "tRACE"},
			wantErr: false,
		},
		{
			name:    "debug",
			args:    args{level: "debug"},
			wantErr: false,
		},
		{
			name:    "info",
			args:    args{level: "info"},
			wantErr: false,
		},
		{
			name:    "warn",
			args:    args{level: "warn"},
			wantErr: false,
		},
		{
			name:    "error",
			args:    args{level: "error"},
			wantErr: false,
		},
		{
			name:    "fatal",
			args:    args{level: "fatal"},
			wantErr: false,
		},
		{
			name:    "panic",
			args:    args{level: "panic"},
			wantErr: false,
		},
		{
			name:    "invalid level",
			args:    args{level: "unknown"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewLogger(tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
