//go:build !integration

package restapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"

	"github.com/Sunaga0709/go-api-template/internal/environment"
	"github.com/Sunaga0709/go-api-template/internal/logger"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	env := mustEnvironment(t, "local")
	lg := testRestAPILogger()

	got := NewServer("8080", env, lg, true, nil, nil)

	if got == nil {
		t.Fatal("NewServer() = nil, want Server")
	}
	if got.router == nil {
		t.Fatal("router = nil, want value")
	}
	if got.port != "8080" {
		t.Fatalf("port = %q, want 8080", got.port)
	}
	if got.environment != env {
		t.Fatalf("environment = %v, want %v", got.environment, env)
	}
	if got.shutdownMW == nil {
		t.Fatal("shutdownMW = nil, want value")
	}
	if !got.router.HideBanner {
		t.Fatal("router.HideBanner = false, want true")
	}
	if !got.router.HidePort {
		t.Fatal("router.HidePort = false, want true")
	}
	if got.router.HTTPErrorHandler == nil {
		t.Fatal("router.HTTPErrorHandler = nil, want value")
	}
}

func TestServer_RoutesHealth(t *testing.T) {
	t.Parallel()

	s := NewServer("8080", mustEnvironment(t, "local"), testRestAPILogger(), true, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got gen.HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got.Status != "ok" {
		t.Fatalf("status body = %q, want ok", got.Status)
	}
}

func TestServer_BeginShutdown(t *testing.T) {
	t.Parallel()

	s := NewServer("8080", mustEnvironment(t, "local"), testRestAPILogger(), true, nil, nil)
	s.BeginShutdown()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}

	var got gen.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if got.Code != http.StatusText(http.StatusServiceUnavailable) {
		t.Fatalf("code = %q, want %q", got.Code, http.StatusText(http.StatusServiceUnavailable))
	}
}

func TestServer_Shutdown(t *testing.T) {
	t.Parallel()

	s := NewServer("8080", mustEnvironment(t, "local"), testRestAPILogger(), true, nil, nil)

	if err := s.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
}

func TestListenAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		port string
		want string
	}{
		{name: "port only", port: "8080", want: ":8080"},
		{name: "trims port", port: " 8080 ", want: ":8080"},
		{name: "host and port", port: "127.0.0.1:8080", want: "127.0.0.1:8080"},
		{name: "ipv6 unspecified", port: "[::]:8080", want: "[::]:8080"},
		{name: "empty", port: "", want: ":"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := listenAddress(tt.port); got != tt.want {
				t.Fatalf("listenAddress(%q) = %q, want %q", tt.port, got, tt.want)
			}
		})
	}
}

func mustEnvironment(t *testing.T, env string) environment.Environment {
	t.Helper()

	e, err := environment.NewEnvironment(env)
	if err != nil {
		t.Fatalf("environment.NewEnvironment() error = %v", err)
	}

	return e
}

func testRestAPILogger() *logger.Logger {
	return &logger.Logger{Inner: zerolog.New(io.Discard)}
}
