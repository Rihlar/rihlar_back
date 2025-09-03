package controllers

import (
	"friend/logger"
	"friend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GenCode(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// 新しいコードを生成する
	code,err := services.GenCode(userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to generate code","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"code": code})
}

func NowCode(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// サービスからコードを取得する
	data, err := services.NowCode(userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get code","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"code": data.Code,"count" : data.UseCount})
}