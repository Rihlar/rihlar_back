package middlewares

import (
	"gcore/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DebugRequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := "userid-79541130-3275-4b90-8677-01323045aca5" //ctx.Request().Header.Get("UserID")

		// トークンを格納
		ctx.Set("token", "")
		// ユーザーIDを格納
		ctx.Set("UserID", userId)

		// 認証処理
		return next(ctx)
	}
}

// 認証ミドルウェア
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// ヘッダからトークンを取得
		token := ctx.Request().Header.Get("Authorization")
		if token == "" {
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// トークンを検証
		claim, err := ValidateToken(token)

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
		}

		// contextにトークンを格納
		ctx.Set("claim", claim)
		// トークンを格納
		ctx.Set("token", token)
		// ユーザーIDを格納
		ctx.Set("UserID", "userid-" + claim.UserID)

		// 認証処理
		return next(ctx)
	}
}