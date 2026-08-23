package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/Sunaga0709/go-api-template/internal/context"
	"github.com/Sunaga0709/go-api-template/internal/restapi/gen"
)

type SetLocationMiddleware struct{}

func NewSetLocationMiddleware() *SetLocationMiddleware {
	return &SetLocationMiddleware{}
}

// SetLocation タイムゾーンを設定する
//
// NOTE:
// サーバ内では基本的にUTCとして扱うため、レスポンスの変換などに使用するタイムゾーン情報をコンテキストにセットする
// ユーザごとにロケーションが異なる場合、認証ミドルウェアと統合も可
// 現状はJST（Asia/Tokyo）で固定
func (s *SetLocationMiddleware) SetLocation(f gen.StrictHandlerFunc, _ string) gen.StrictHandlerFunc {
	return func(ctx echo.Context, request any) (any, error) {
		loc, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			loc = time.FixedZone("JST", 9*60*60)
		}

		setLocationCtx := context.SetLocation(ctx.Request().Context(), loc)
		ctx.SetRequest(ctx.Request().WithContext(setLocationCtx))

		return f(ctx, request)
	}
}
