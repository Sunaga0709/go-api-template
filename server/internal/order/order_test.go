//go:build !integration

package order

import (
	"errors"
	"fmt"
	"testing"
)

func TestOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		o    Order
		want int
	}{
		// 正常系
		{
			name: "asc",
			o:    OrderAsc,
			want: 0,
		},
		{
			name: "desc",
			o:    OrderDesc,
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := int(tt.o); got != tt.want {
				t.Errorf("Order = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOrder(t *testing.T) {
	t.Parallel()

	type args struct {
		order string
	}
	tests := []struct {
		name          string
		args          args
		want          Order
		wantErr       bool
		wantErrString string
	}{
		// 正常系
		{
			name:    "asc",
			args:    args{order: orderAscString},
			want:    OrderAsc,
			wantErr: false,
		},
		{
			name:    "desc",
			args:    args{order: orderDescString},
			want:    OrderDesc,
			wantErr: false,
		},
		// 異常系
		{
			name:          "empty string",
			args:          args{order: ""},
			want:          0,
			wantErr:       true,
			wantErrString: "order error: invalid order string: got = ",
		},
		{
			name:          "unknown string",
			args:          args{order: "unknown"},
			want:          0,
			wantErr:       true,
			wantErrString: "order error: invalid order string: got = unknown",
		},
		{
			name:          "case sensitive (Asc)",
			args:          args{order: "Asc"},
			want:          0,
			wantErr:       true,
			wantErrString: "order error: invalid order string: got = Asc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseOrder(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOrder() = %v, want %v", got, tt.want)
			}
			if tt.wantErrString != "" {
				if err == nil {
					t.Fatalf("ParseOrder() error = nil, want %v", tt.wantErrString)
				}
				if err.Error() != tt.wantErrString {
					t.Errorf("ParseOrder() error = %v, want %v", err, tt.wantErrString)
				}
				var orderErr Error
				if !errors.As(err, &orderErr) {
					t.Errorf("ParseOrder() error type = %T, want Error", err)
				}
			}
		})
	}
}

func Test_newError(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := newError(tt.args.err)
			if !errors.Is(got.error, tt.want.error) {
				t.Errorf("newError() error = %v, want %v", got.error, tt.want.error)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    Error
		want string
	}{
		// 正常系
		{
			name: "error",
			e:    Error{error: errors.New("base error")},
			want: "order error: base error",
		},
		{
			name: "wrapped error",
			e:    Error{error: fmt.Errorf("outer: %w", errors.New("inner"))},
			want: "order error: outer: inner",
		},
		// 境界値（nil）
		{
			name: "nil error",
			e:    Error{error: nil},
			want: "order error: <nil>",
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
