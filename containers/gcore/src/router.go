package main

import (
	"gcore/controllers"

	"github.com/labstack/echo/v4"
)

func InitRouter(router *echo.Echo) *echo.Echo {
	// 位置情報が更新されるたびに呼ばれるエンドポイント
	router.POST("/report/movement", controllers.ReportMovement)

	return router
}