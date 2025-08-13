package middlewares

import (
	"game/logger"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

func DebugRequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userId := "userid-79541130-3275-4b90-8677-01323045aca5" //ctx.Request().Header.Get("UserID")

		// トークンを格納
		ctx.Set("token", "")
		// ユーザーIDを格納
		ctx.Set("UserID", userId)

		// コンテキスト登録処理
		circleID := ctx.Request().Header.Get("CircleID")
		ctx.Set("CircleID", circleID)

		gameID := ctx.Request().Header.Get("GameID")
		ctx.Set("GameID", gameID)

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


		// コンテキスト登録処理
		circleID := ctx.Request().Header.Get("CircleID")
		ctx.Set("CircleID", circleID)

		gameID := ctx.Request().Header.Get("GameID")
		ctx.Set("GameID", gameID)

		// 認証処理
		return next(ctx)
	}
}

func RequireLabel(labels []string) echo.MiddlewareFunc {
	return RequireLabelMiddleware{
		Labels: labels,
	}.RequireAuth
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

		if len(claim.Labels) == 0 {
			return ctx.JSON(http.StatusForbidden, echo.Map{"error": "you don't have permission"})
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

		// コンテキスト登録処理
		circleID := ctx.Request().Header.Get("CircleID")
		ctx.Set("CircleID", circleID)

		gameID := ctx.Request().Header.Get("GameID")
		ctx.Set("GameID", gameID)

		// 認証処理
		return next(ctx)
	}
}
