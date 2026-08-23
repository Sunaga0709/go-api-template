//go:build !integration

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	appcontext "github.com/Sunaga0709/go-api-template/internal/context"
)

func TestSetLocationMiddleware_SetLocation(t *testing.T) {
	t.Parallel()

	mw := NewSetLocationMiddleware()
	handler := mw.SetLocation(func(ctx echo.Context, request any) (any, error) {
		loc := appcontext.GetLocation(ctx.Request().Context())
		if loc == nil {
			t.Fatal("location = nil, want Asia/Tokyo")
		}
		if loc.String() != "Asia/Tokyo" {
			t.Fatalf("location = %s, want Asia/Tokyo", loc)
		}

		return request, nil
	}, "operation")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := echo.New().NewContext(req, httptest.NewRecorder())

	got, err := handler(ctx, "ok")
	if err != nil {
		t.Fatalf("SetLocation() error = %v", err)
	}
	if got != "ok" {
		t.Fatalf("SetLocation() response = %v, want ok", got)
	}
	if appcontext.GetLocation(req.Context()) != nil {
		t.Fatal("original request context location was mutated")
	}
}
