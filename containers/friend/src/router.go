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

	// 現在のコードを取得する関数
	router.GET("/nowcode",controllers.NowCode)

	// フレンドリクエストを送信する関数
	router.POST("/request",controllers.SendRequest)

	// フレンドリストの取得
	router.GET("/list",controllers.GetFriendList)

	// フレンドの削除
	router.DELETE("/delete",controllers.DeleteFriend)

	return router
}
