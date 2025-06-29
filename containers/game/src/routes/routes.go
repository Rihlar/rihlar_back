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
	// rankingTop10
	router.GET("/ranking/top/:game_id", controllers.GetRankingTopHandler)
	
	//円詳細取得
	router.GET("/circle/:circle_id", controllers.GetCircleDeteileHandler)

	return router
}
