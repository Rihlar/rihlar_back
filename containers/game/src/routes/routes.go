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
	//円詳細取得
	router.GET("/circle/:circle_id", controllers.GetCircleDeteileHandler)


	//adminゲーム作成
	router.POST("game/create", controllers.CreateAdminGameHandler)
	// 終了済みゲーム一覧
	router.GET("/endgame/:user_id", controllers.GetEndGamesHandler)
	// 参加している全てのゲーム一覧取得
	router.GET("/joingame/:user_id", controllers.GetJoinGamesHandler)


	return router
}
