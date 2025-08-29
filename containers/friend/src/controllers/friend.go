package controllers

import (
	"friend/logger"
	"friend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// フレンドリクエストを送信する関数
func SendRequest(ctx echo.Context) error {
	// コードを取得する
	friendCode := ctx.Request().Header.Get("code")

	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	err := services.SendRequest(userId, friendCode)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to send request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"result": "success"})
}

// フレンドのリストを取得する関数
func GetFriendList(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	data, err := services.GetFriendList(userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get friend list","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}