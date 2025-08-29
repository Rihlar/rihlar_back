package controllers

import (
	"friend/logger"
	"friend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 受信済みリクエストを取得する関数
func GetRecvedRequest(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	data, err := services.GetRecvedRequest(userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get recved request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}