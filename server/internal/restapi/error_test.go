//go:build !integration

package restapi

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
			name: "simple error",
			e:    Error{error: errors.New("base error")},
			want: "restapi error: base error",
		},
		{
			name: "wrapped error",
			e:    Error{error: fmt.Errorf("outer: %w", errors.New("inner"))},
			want: "restapi error: outer: inner",
		},
		{
			name: "nil error",
			e:    Error{error: nil},
			want: "restapi error: <nil>",
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
