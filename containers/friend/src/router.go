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

	// リクエスト用グループ作成
	requestg := router.Group("/request")
	{
		// 受信済みフレンドリクエスト取得
		requestg.GET("/recved",controllers.GetRecvedRequest)

		// フレンドリクエストを拒否する関数
		requestg.POST("/reject",controllers.RejectRequest)

		// フレンドリクエストを承認する関数
		requestg.POST("/accept",controllers.AcceptRequest)
	}

	return router
}
