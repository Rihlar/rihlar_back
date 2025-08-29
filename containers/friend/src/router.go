package main

import (
	"friend/controllers"
	"friend/middlewares"

	"github.com/labstack/echo/v4"
)

func InitRoutes(router *echo.Echo) *echo.Echo {
	// 認証を必要とするように変更
	router.Use(middlewares.RequireAuth)

	// コードを生成するエンドポイント
	router.GET("/gencode",controllers.GenCode)

	return router
}
