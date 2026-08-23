//go:build !integration

package controller

import (
	"context"
	"testing"

	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func TestController_GetHealth(t *testing.T) {
	t.Parallel()

	c := NewController(nil, nil, nil, nil, nil)

	res, err := c.GetHealth(context.Background(), gen.GetHealthRequestObject{})
	if err != nil {
		t.Fatalf("GetHealth() error = %v", err)
	}

	got, ok := res.(gen.GetHealth200JSONResponse)
	if !ok {
		t.Fatalf("GetHealth() response type = %T, want %T", res, gen.GetHealth200JSONResponse{})
	}
	if got.Status != "ok" {
		t.Fatalf("GetHealth() status = %q, want ok", got.Status)
	}
}
