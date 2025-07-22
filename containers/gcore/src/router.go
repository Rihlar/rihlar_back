package main

import (
	"gcore/controllers"
	// "gcore/middlewares"

	"github.com/labstack/echo/v4"
)

func InitRouter(router *echo.Echo) *echo.Echo {
	// ここにルーティング関連を書く

	// 認証を必要とする
	// router.Use(middlewares.RequireAuth)

	// 移動を報告するエンドポイント
	router.POST("/report/movement", controllers.ReportMovement)

	// 移動を取得するエンドポイント
	router.GET("/get/movement", controllers.GetReportedMovement)

	// 円を作成するエンドポイント
	router.POST("/create/circle", controllers.CreateCircle)

	// TOP10を取得するエンドポイント
	router.GET("/get/top", controllers.GetRanking)

	return router
}
