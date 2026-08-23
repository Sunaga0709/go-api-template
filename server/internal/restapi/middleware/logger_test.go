//go:build !integration

package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestLoggerMiddleware_AccessLogging(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		hiddenSuccess bool
		handlerErr    error
		wantResponse  any
		wantLogged    bool
		wantLevel     string
		wantStatus    string
	}{
		{
			name:         "logs success",
			wantResponse: "ok",
			wantLogged:   true,
			wantLevel:    "info",
			wantStatus:   `"status":200`,
		},
		{
			name:          "hides success",
			hiddenSuccess: true,
			wantResponse:  "ok",
			wantLogged:    false,
		},
		{
			name:       "logs error detail as warning",
			handlerErr: NewErrorDetail(http.StatusBadRequest, errors.New("bad request")),
			wantLogged: true,
			wantLevel:  "warn",
			wantStatus: `"status":400`,
		},
		{
			name:       "logs plain error as error",
			handlerErr: errors.New("failed"),
			wantLogged: true,
			wantLevel:  "error",
			wantStatus: `"status":500`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			mw := NewLoggerMiddleware(testLogger(t, &buf), tt.hiddenSuccess)
			handler := mw.AccessLogging(func(_ echo.Context, _ any) (any, error) {
				if tt.handlerErr != nil {
					return nil, tt.handlerErr
				}
				return tt.wantResponse, nil
			}, "operation-id")

			ctx := echo.New().NewContext(
				httptest.NewRequest(http.MethodPost, "/v1/test", nil),
				httptest.NewRecorder(),
			)
			got, err := handler(ctx, nil)

			if !errors.Is(err, tt.handlerErr) {
				t.Fatalf("AccessLogging() error = %v, want %v", err, tt.handlerErr)
			}
			if got != tt.wantResponse {
				t.Fatalf("AccessLogging() response = %v, want %v", got, tt.wantResponse)
			}

			logged := buf.String()
			if tt.wantLogged && logged == "" {
				t.Fatal("log is empty, want output")
			}
			if !tt.wantLogged {
				if logged != "" {
					t.Fatalf("log = %q, want empty", logged)
				}
				return
			}

			wants := []string{
				`"level":"` + tt.wantLevel + `"`,
				`"operation_id":"operation-id"`,
				`"method":"POST"`,
				`"path":"/v1/test"`,
				tt.wantStatus,
				`"timestamp":`,
				`"elapsed":`,
			}
			for _, want := range wants {
				if !strings.Contains(logged, want) {
					t.Errorf("log does not contain %q: %s", want, logged)
				}
			}
		})
	}
}

func TestLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status int
		want   string
	}{
		{status: http.StatusOK, want: "info"},
		{status: http.StatusBadRequest, want: "warn"},
		{status: http.StatusInternalServerError, want: "error"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			if got := logLevel(tt.status).String(); got != tt.want {
				t.Fatalf("logLevel(%d) = %s, want %s", tt.status, got, tt.want)
			}
		})
	}
}
