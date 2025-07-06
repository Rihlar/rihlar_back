package routes

// エンドポイントのルーティング
import (
	"game/controllers"
	"game/middlewares"

	"github.com/labstack/echo/v4"
)

// ルーティング
func InitRoutes(router *echo.Echo) *echo.Echo {

	//ranking取得
	router.GET("/ranking/personal/:user_id", controllers.GetMyRankingHandler)
	
	// 上位３位の円を取得する
	router.GET("/ranking/top/:game_id", controllers.GetRankingTopHandler)

	// rankingTop10
	// router.GET("/ranking/top/:game_id", controllers.GetRankingTopHandler)

	// ゲームを作成するエンドポイント
	router.POST("/create", controllers.CreateGameHandler,middlewares.RequireLabel([]string{"admin"}))

	// ゲームに参加するエンドポイント
	router.POST("/join", controllers.JoinGameHandler)
	
	//円詳細取得
	router.GET("/circle/:circle_id", controllers.GetCircleDeteileHandler)

	// 終了済みゲーム一覧
	router.GET("/endgame/:user_id", controllers.GetEndGamesHandler)
	// 参加している全てのゲーム一覧取得
	router.GET("/joingame/:user_id", controllers.GetJoinGamesHandler)

	return router
}
