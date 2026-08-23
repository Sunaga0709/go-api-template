//go:build !integration

package main

import (
	"errors"
	"fmt"
	"testing"
)

func Test_newError(t *testing.T) {
	baseErr := errors.New("base error")
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		// 正常系
		{
			name: "wraps error",
			args: args{err: baseErr},
			want: Error{error: baseErr},
		},
		// 境界値（nil）
		{
			name: "nil error",
			args: args{err: nil},
			want: Error{error: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newError(tt.args.err); !errors.Is(got.error, tt.want.error) {
				t.Errorf("newError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    Error
		want string
	}{
		// 正常系
		{
			name: "simple error",
			e:    Error{error: errors.New("base error")},
			want: "server error: base error",
		},
		{
			name: "wrapped error",
			e:    Error{error: fmt.Errorf("outer: %w", errors.New("inner"))},
			want: "server error: outer: inner",
		},
		// 境界値（nil）
		{
			name: "nil error",
			e:    Error{error: nil},
			want: "server error: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
