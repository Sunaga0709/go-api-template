//go:build !integration

package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/Sunaga0709/go-api-template/internal/logger"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func TestPanicRecoveryMiddleware_Recover(t *testing.T) {
	tests := []struct {
		name       string
		handler    echo.HandlerFunc
		wantCalled bool
		wantErr    bool
		wantStatus int
	}{
		{
			name: "passes request without panic",
			handler: func(_ echo.Context) error {
				return nil
			},
			wantCalled: true,
			wantErr:    false,
		},
		{
			name: "converts panic to internal server error",
			handler: func(_ echo.Context) error {
				panic("unexpected panic")
			},
			wantCalled: true,
			wantErr:    true,
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := NewPanicRecoveryMiddleware(testLogger(t, nil))

			called := false
			handler := mw.Recover(func(ctx echo.Context) error {
				called = true
				return tt.handler(ctx)
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
				t.Fatalf("Recover() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				return
			}

			var detail ErrorDetail
			if !errors.As(err, &detail) {
				t.Fatalf("Recover() error = %T, want ErrorDetail", err)
			}
			if detail.Status != tt.wantStatus {
				t.Errorf("Recover() status = %v, want %v", detail.Status, tt.wantStatus)
			}
		})
	}
}

func TestPanicRecoveryMiddleware_RecoverWritesInternalServerError(t *testing.T) {
	e := echo.New()
	errorHandling := NewErrorHandlingMiddleware()
	panicRecovery := NewPanicRecoveryMiddleware(testLogger(t, nil))

	e.HTTPErrorHandler = errorHandling.Handle
	e.Use(panicRecovery.Recover)
	e.GET("/panic", func(_ echo.Context) error {
		panic("unexpected panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %v, want %v", rec.Code, http.StatusInternalServerError)
	}

	var got gen.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if got.Code != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("code = %v, want %v", got.Code, http.StatusText(http.StatusInternalServerError))
	}
	if got.Message != defaultClientErrMessage500 {
		t.Errorf("message = %v, want %v", got.Message, defaultClientErrMessage500)
	}
}

func TestPanicRecoveryMiddleware_RecoverLogsPanic(t *testing.T) {
	var buf bytes.Buffer
	mw := NewPanicRecoveryMiddleware(testLogger(t, &buf))

	handler := mw.Recover(func(_ echo.Context) error {
		panic("unexpected panic")
	})
	ctx := echo.New().NewContext(
		httptest.NewRequest(http.MethodPost, "/panic", nil),
		httptest.NewRecorder(),
	)

	if err := handler(ctx); err == nil {
		t.Fatal("Recover() error = nil, want error")
	}

	got := buf.String()
	wants := []string{
		`"level":"error"`,
		`"event":"panic_recovered"`,
		`"method":"POST"`,
		`"path":"/panic"`,
		`"status":500`,
		`"panic":"unexpected panic"`,
		`"stack":`,
	}
	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf("log does not contain %q: %s", want, got)
		}
	}
}

func testLogger(t *testing.T, buf *bytes.Buffer) logger.Logger {
	t.Helper()

	if buf == nil {
		buf = &bytes.Buffer{}
	}

	return logger.Logger{
		Inner: zerolog.New(buf).Level(zerolog.TraceLevel),
	}
}
