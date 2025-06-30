package routes

import (
	"user/controllers"

	"github.com/labstack/echo/v4"
)

func InitRoute(router *echo.Echo) *echo.Echo {
	//profile全取得
	router.GET("/profile/:user_id", controllers.GetProfileById)

	router.PUT("/profile/:user_id", controllers.UpdateProfileById)

	return router
}
