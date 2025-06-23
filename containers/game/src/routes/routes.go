package routes

// エンドポイントのルーティング
import (
	"game/controllers"
	"github.com/labstack/echo/v4"
)

// ルーティング　
func InitRoutes(router *echo.Echo) *echo.Echo {

	//ranking取得
	router.GET("/ranking/personal/:user_id", controllers.GetMyRankingHandler)

	return router
}
