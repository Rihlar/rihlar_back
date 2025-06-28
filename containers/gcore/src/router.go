package main

import (
	"gcore/controllers"

	"github.com/labstack/echo/v4"
)

func InitRouter(router *echo.Echo) *echo.Echo {
	// ここにルーティング関連を書く

	// 移動を報告するエンドポイント
	router.POST("/report/movement", controllers.ReportMovement)

	return router
}