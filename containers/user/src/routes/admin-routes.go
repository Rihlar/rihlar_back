package routes

import (
	"user/controllers"
	"user/middlewares"

	"github.com/labstack/echo/v4"
)

func InitAdminRoute(router *echo.Echo) *echo.Echo {
	// グループを作成する
	adminG := router.Group("/admin")
	{
		// 認証を必要とする
		adminG.Use(middlewares.RequireLabel([]string{"admin"}))

				//profile一件取得
				adminG.GET("/profiles", controllers.GetAllProfilesFromAdmin)
		
				// 全てのユーザー情報を取得する
				adminG.GET("/users", controllers.GetAllUsersFromAdmin)
		
				//profile作成		adminG.POST("/profile", controllers.CreateProfileFromAdmin)

		// profile削除
		adminG.DELETE("/profile", controllers.DeleteProfileFromAdmin)

		// リージョン一覧を取得
		adminG.GET("/regions", controllers.GetRegionList)
	}

	return router
}
