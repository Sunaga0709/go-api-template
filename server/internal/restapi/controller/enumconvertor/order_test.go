//go:build !integration

package enumconvertor

import (
	"testing"

	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func TestConvertOrderToPrimitive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ord  gen.Order
		want string
	}{
		{name: "asc", ord: gen.OrderAsc, want: "asc"},
		{name: "desc", ord: gen.OrderDesc, want: "desc"},
		{name: "unknown", ord: gen.Order("unknown"), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ConvertOrderToPrimitive(tt.ord); got != tt.want {
				t.Fatalf("ConvertOrderToPrimitive() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConvertOrderToOpenAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ord  string
		want gen.Order
	}{
		{name: "desc", ord: "desc", want: gen.OrderDesc},
		{name: "asc", ord: "asc", want: gen.OrderAsc},
		{name: "unknown defaults to asc", ord: "unknown", want: gen.OrderAsc},
		{name: "empty defaults to asc", ord: "", want: gen.OrderAsc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ConvertOrderToOpenAPI(tt.ord); got != tt.want {
				t.Fatalf("ConvertOrderToOpenAPI() = %q, want %q", got, tt.want)
			}
		})
	}
}
