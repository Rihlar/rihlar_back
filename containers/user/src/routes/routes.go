package routes

import (
	"user/controllers"
	"user/middlewares"

	"github.com/labstack/echo/v4"
)

func InitRoute(router *echo.Echo) *echo.Echo {
	// 認証を必要とする
	router.Use(middlewares.RequireAuth)

	//profile一件取得
	router.GET("/profile", controllers.GetProfileById)
	//profile作成
	router.POST("/profile", controllers.CreateProfile)
	//profile編集
	router.PUT("/profile", controllers.UpdateProfileById)
	//privacyエリア取得
	router.GET("/privacy", controllers.GetPrivacyProfile)
	//privacyエリア編集
	router.PUT("/privacy", controllers.UpdatePrivacyProfile)
	//所属地域の取得
	router.GET("/region", controllers.GetRegionProfile)
	//所属地域の編集
	router.PUT("/region", controllers.UpdateRegionProfile)
	//実績の取得
	router.GET("/achive", controllers.GetAchiveProfile)
	//実績の編集
	router.PUT("/achive", controllers.UpdateAchiveProfile)

	return router
}
