//go:build !integration

package repository

import (
	"errors"
	"fmt"
	"testing"
)

func TestError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Error
		want string
	}{
		{
			name: "unknown",
			e:    Error{Kind: ErrorKindUnknown, error: errors.New("base error")},
			want: "query error (kind=unknown): base error",
		},
		{
			name: "get",
			e:    Error{Kind: ErrorKindGet, error: errors.New("base error")},
			want: "query error (kind=get): base error",
		},
		{
			name: "not_found",
			e:    Error{Kind: ErrorKindNotFound, error: errors.New("base error")},
			want: "query error (kind=not_found): base error",
		},
		{
			name: "convert_domain_model",
			e:    Error{Kind: ErrorKindConvertDomainModel, error: errors.New("base error")},
			want: "query error (kind=convert_domain_model): base error",
		},
		{
			name: "wrapped error",
			e:    Error{Kind: ErrorKindGet, error: fmt.Errorf("outer: %w", errors.New("inner"))},
			want: "query error (kind=get): outer: inner",
		},
		{
			name: "nil error",
			e:    Error{Kind: ErrorKindUnknown, error: nil},
			want: "query error (kind=unknown): <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base error")
	got := NewError(baseErr, ErrorKindGet)

	if got.Kind != ErrorKindGet {
		t.Errorf("NewError() Kind = %v, want %v", got.Kind, ErrorKindGet)
	}
	if got.Error() != "query error (kind=get): base error" {
		t.Errorf("NewError().Error() = %v, want query error (kind=get): base error", got.Error())
	}
}

func TestErrorConstructors(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base error")
	tests := []struct {
		name string
		fn   func(error) Error
		want ErrorKind
	}{
		{
			name: "unknown",
			fn:   NewUnknownError,
			want: ErrorKindUnknown,
		},
		{
			name: "get",
			fn:   NewGetError,
			want: ErrorKindGet,
		},
		{
			name: "not found",
			fn:   NewNotFoundError,
			want: ErrorKindNotFound,
		},
		{
			name: "convert domain model",
			fn:   NewConvertDomainModelError,
			want: ErrorKindConvertDomainModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.fn(baseErr)
			if got.Kind != tt.want {
				t.Errorf("%s Kind = %v, want %v", tt.name, got.Kind, tt.want)
			}
			if got.Error() != fmt.Sprintf("query error (kind=%s): base error", tt.want.String()) {
				t.Errorf("%s Error() = %v", tt.name, got.Error())
			}
		})
	}
}
