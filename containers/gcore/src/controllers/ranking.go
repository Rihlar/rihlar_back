package controllers

import (
	"gcore/logger"
	"gcore/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRanking(ctx echo.Context) error {
	// ユーザーID取得
	userID := ctx.Request().Header.Get("UserID")

	// ゲームID取得
	gameID := ctx.Request().Header.Get("GameID")

	// サービスを呼び出す
	ranking, err := services.GetRanking(gameID, userID)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"result": "success",
		"data": ranking,
	})
}