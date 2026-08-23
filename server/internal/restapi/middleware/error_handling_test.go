//go:build !integration

package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

func TestErrorDetail_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ErrorDetail
		want string
	}{
		{
			name: "bad request",
			e:    ErrorDetail{Status: http.StatusBadRequest, error: errors.New("base error")},
			want: "error detail (status=400): base error",
		},
		{
			name: "internal server error",
			e:    ErrorDetail{Status: http.StatusInternalServerError, error: errors.New("base error")},
			want: "error detail (status=500): base error",
		},
		{
			name: "wrapped error",
			e:    ErrorDetail{Status: http.StatusInternalServerError, error: fmt.Errorf("outer: %w", errors.New("inner"))},
			want: "error detail (status=500): outer: inner",
		},
		{
			name: "nil error",
			e:    ErrorDetail{Status: http.StatusInternalServerError, error: nil},
			want: "error detail (status=500): <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.e.Error(); got != tt.want {
				t.Errorf("ErrorDetail.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorHandlingMiddleware_ConvertOpenAPI(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("failed")
	tests := []struct {
		name       string
		handlerErr error
		wantRes    any
		wantStatus int
	}{
		{
			name:    "returns response without error",
			wantRes: "ok",
		},
		{
			name:       "keeps error detail",
			handlerErr: NewErrorDetail(http.StatusBadRequest, baseErr),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "converts plain error to internal server error",
			handlerErr: baseErr,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mw := NewErrorHandlingMiddleware()
			handler := mw.ConvertOpenAPI(func(_ echo.Context, _ any) (any, error) {
				if tt.handlerErr != nil {
					return nil, tt.handlerErr
				}
				return tt.wantRes, nil
			}, "operation")

			got, err := handler(echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder()), nil)
			if tt.handlerErr == nil {
				if err != nil {
					t.Fatalf("ConvertOpenAPI() error = %v", err)
				}
				if got != tt.wantRes {
					t.Fatalf("ConvertOpenAPI() response = %v, want %v", got, tt.wantRes)
				}
				return
			}

			if got != nil {
				t.Fatalf("ConvertOpenAPI() response = %#v, want nil", got)
			}
			assertErrorDetailStatus(t, err, tt.wantStatus)
		})
	}
}

func TestErrorHandlingMiddleware_Handle(t *testing.T) {
	t.Parallel()

	customMessage := "custom message"
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{
			name:        "bad request",
			err:         NewErrorDetail(http.StatusBadRequest, errors.New("bad request")),
			wantStatus:  http.StatusBadRequest,
			wantCode:    http.StatusText(http.StatusBadRequest),
			wantMessage: defaultClientErrMessage400,
		},
		{
			name:        "not found",
			err:         NewErrorDetail(http.StatusNotFound, errors.New("not found")),
			wantStatus:  http.StatusNotFound,
			wantCode:    http.StatusText(http.StatusNotFound),
			wantMessage: defaultClientErrMessage404,
		},
		{
			name:        "conflict with client message",
			err:         NewErrorDetail(http.StatusConflict, errors.New("conflict")).WithClientMessage(customMessage),
			wantStatus:  http.StatusConflict,
			wantCode:    http.StatusText(http.StatusConflict),
			wantMessage: customMessage,
		},
		{
			name:        "service unavailable",
			err:         NewErrorDetail(http.StatusServiceUnavailable, errors.New("shutdown")),
			wantStatus:  http.StatusServiceUnavailable,
			wantCode:    http.StatusText(http.StatusServiceUnavailable),
			wantMessage: defaultClientErrMessage503,
		},
		{
			name:        "plain error",
			err:         errors.New("failed"),
			wantStatus:  http.StatusInternalServerError,
			wantCode:    http.StatusText(http.StatusInternalServerError),
			wantMessage: defaultClientErrMessage500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			rec := httptest.NewRecorder()
			ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)

			NewErrorHandlingMiddleware().Handle(tt.err, ctx)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var got gen.ErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", got.Code, tt.wantCode)
			}
			if got.Message != tt.wantMessage {
				t.Errorf("message = %q, want %q", got.Message, tt.wantMessage)
			}
		})
	}
}

func TestErrorHandlingMiddleware_HandleCommittedResponse(t *testing.T) {
	t.Parallel()

	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
	ctx.Response().Committed = true

	NewErrorHandlingMiddleware().Handle(NewErrorDetail(http.StatusBadRequest, errors.New("failed")), ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want untouched recorder default %d", rec.Code, http.StatusOK)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", rec.Body.String())
	}
}

func TestErrorDetail_Unwrap(t *testing.T) {
	t.Parallel()

	baseErr := errors.New("base")
	detail := NewErrorDetail(http.StatusInternalServerError, baseErr)

	if !errors.Is(detail, baseErr) {
		t.Fatalf("errors.Is(detail, baseErr) = false, want true")
	}
}

func assertErrorDetailStatus(t *testing.T, err error, want int) {
	t.Helper()

	var detail ErrorDetail
	if !errors.As(err, &detail) {
		t.Fatalf("error type = %T, want ErrorDetail", err)
	}
	if detail.Status != want {
		t.Fatalf("status = %d, want %d", detail.Status, want)
	}
}
