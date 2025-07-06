package middlewares

import (
	"game/logger"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

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
		ctx.Set("UserID", claim.UserID)

		// 認証処理
		return next(ctx)
	}
}

func RequireLabel(labels []string) echo.MiddlewareFunc {
	return RequireLabelMiddleware{labels}.RequireAuth
}

type RequireLabelMiddleware struct {
	Labels []string
}

// 認証ミドルウェア
func (middleware RequireLabelMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
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

		// 特定のラベルを持っているか検証
		for _, label := range claim.Labels {
			// 対応したラベルに含まれているか
			if !slices.Contains(middleware.Labels, label) {
				// 含まれない場合
				return ctx.JSON(http.StatusForbidden, echo.Map{"error": "you don't have permission"})
			}
		}

		// contextにトークンを格納
		ctx.Set("claim", claim)
		// トークンを格納
		ctx.Set("token", token)
		// ユーザーIDを格納
		ctx.Set("UserID", claim.UserID)

		// 認証処理
		return next(ctx)
	}
}