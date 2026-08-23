//go:build !integration

package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestShutdownMiddleware_RejectDuringShutdown(t *testing.T) {
	type fields struct {
		beginShutdown bool
	}
	tests := []struct {
		name       string
		fields     fields
		wantCalled bool
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "passes request before shutdown",
			fields:     fields{beginShutdown: false},
			wantCalled: true,
			wantErr:    false,
		},
		{
			name:       "rejects request during shutdown",
			fields:     fields{beginShutdown: true},
			wantCalled: false,
			wantErr:    true,
			wantStatus: http.StatusServiceUnavailable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := NewShutdownMiddleware()
			if tt.fields.beginShutdown {
				mw.BeginShutdown()
			}

			called := false
			handler := mw.RejectDuringShutdown(func(_ echo.Context) error {
				called = true
				return nil
			})

			ctx := echo.New().NewContext(
				httptest.NewRequest(http.MethodGet, "/health", nil),
				httptest.NewRecorder(),
			)
			err := handler(ctx)

			if called != tt.wantCalled {
				t.Errorf("handler called = %v, want %v", called, tt.wantCalled)
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("RejectDuringShutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				return
			}

			var detail ErrorDetail
			if !errors.As(err, &detail) {
				t.Fatalf("RejectDuringShutdown() error = %T, want ErrorDetail", err)
			}
			if detail.Status != tt.wantStatus {
				t.Errorf("RejectDuringShutdown() status = %v, want %v", detail.Status, tt.wantStatus)
			}
		})
	}
}
