package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

const (
	defaultClientErrMessage400 = "不正な値がリクエストされました。"
	defaultClientErrMessage404 = "リソースが見つかりませんでした。"
	defaultClientErrMessage409 = "既にリソースが存在しています。"
	defaultClientErrMessage503 = "サーバ停止中です。"
	defaultClientErrMessage500 = "サーバエラーが発生しました。"
)

type ErrorHandlingMiddleware struct{}

func NewErrorHandlingMiddleware() *ErrorHandlingMiddleware {
	return &ErrorHandlingMiddleware{}
}

// ConvertOpenAPI ハンドラー（コントローラー）が返すエラーを、必ずステータスを持つ ErrorDetail に正規化する
// ここではHTTPレスポンスを書き込まず（書き込みは Handle が担う）
// 後段のミドルウェアがステータスを参照できるようにすることに専念する
func (c *ErrorHandlingMiddleware) ConvertOpenAPI(f gen.StrictHandlerFunc, _ string) gen.StrictHandlerFunc {
	return func(ctx echo.Context, request any) (any, error) {
		res, err := f(ctx, request)
		if err == nil {
			return res, nil
		}

		if detail, ok := errors.AsType[ErrorDetail](err); ok {
			return nil, detail
		}

		return nil, NewErrorDetail(http.StatusInternalServerError, err)
	}
}

// Handle echoのHTTPErrorHandlerとして登録し、OpenAPI定義の ErrorResponse を書き出す
func (c *ErrorHandlingMiddleware) Handle(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	var detail ErrorDetail
	if !errors.As(err, &detail) {
		detail = NewErrorDetail(http.StatusInternalServerError, err)
	}

	var clientMessage string
	switch detail.Status {
	case http.StatusBadRequest:
		clientMessage = defaultClientErrMessage400
	case http.StatusNotFound:
		clientMessage = defaultClientErrMessage404
	case http.StatusConflict:
		clientMessage = defaultClientErrMessage409
	case http.StatusServiceUnavailable:
		clientMessage = defaultClientErrMessage503
	default:
		clientMessage = defaultClientErrMessage500
	}

	if detail.ClientMessage != nil {
		clientMessage = *detail.ClientMessage
	}

	res := gen.ErrorResponse{
		Code:    http.StatusText(detail.Status),
		Message: clientMessage,
	}

	ctx.JSON(detail.Status, res) //nolint:errcheck // エラーハンドラの最終出力のためエラーを返しても処理できない
}

//nolint:errname // エラー詳細を表す構造体でありXxxError命名は不適切なため
type ErrorDetail struct {
	Status        int
	error         error
	ClientMessage *string
}

func NewErrorDetail(status int, err error) ErrorDetail {
	return ErrorDetail{
		Status:        status,
		error:         err,
		ClientMessage: nil,
	}
}

func (e ErrorDetail) WithClientMessage(msg string) ErrorDetail {
	err := e
	err.ClientMessage = &msg
	return err
}

func (e ErrorDetail) Error() string {
	return fmt.Sprintf("error detail (status=%d): %v", e.Status, e.error)
}

func (e ErrorDetail) Unwrap() error {
	return e.error
}
