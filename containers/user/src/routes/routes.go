package routes

import (
	"user/controllers"

	"github.com/labstack/echo/v4"
)

func InitRoute(router *echo.Echo) *echo.Echo {
	//profile全取得
	router.GET("/profile/:user_id", controllers.GetProfileById)
	//profile編集
	router.PUT("/profile/:user_id", controllers.UpdateProfileById)
	//privacyエリア取得
	router.GET("/privacy/:user_id", controllers.GetPrivacyProfile)
	//privacyエリア編集
	router.PUT("/privacy/:user_id", controllers.UpdatePrivacyProfile)
	//所属地域の取得
	router.GET("region/:user_id", controllers.GetRegionProfile)
	//所属地域の編集
	router.PUT("region/:user_id", controllers.UpdateRegionProfile)

	return router
}
