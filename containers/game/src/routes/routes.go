package routes

// エンドポイントのルーティング
import (
	"game/controllers"
	"game/middlewares"

	"github.com/labstack/echo/v4"
)

// ルーティング
func InitRoutes(router *echo.Echo) *echo.Echo {
	// 認証を必要とするように変更
	// router.Use(middlewares.RequireAuth)

	//ranking取得
	router.GET("/ranking/personal/:user_id", controllers.GetMyRankingHandler)

	// 上位３位の円を取得する
	router.GET("/ranking/top", controllers.GetRankingTopHandler, middlewares.DebugRequireAuth)

	// rankingTop10
	router.GET("/ranking/top10", controllers.GetRankingTop10Handler, middlewares.RequireAuth)

	// ソロ用のランキングtop10
	router.GET("/ranking/solo/top10", controllers.GetRankingTop10SoloHandler, middlewares.RequireAuth)

	// ゲームの一覧を取得するエンドポイント (モバイル用)
	router.GET("/allgames", controllers.GetAllGameListHandler, middlewares.RequireAuth)

	// ゲームを作成するエンドポイント
	router.POST("/create", controllers.CreateGameHandler, middlewares.RequireLabel([]string{"admin"}))

	// ゲームに参加するエンドポイント
	router.POST("/join", controllers.JoinGameHandler, middlewares.RequireAuth)

	// ゲームの一覧を取得するエンドポイント
	router.GET("/list", controllers.GetGameListHandler, middlewares.RequireLabel([]string{"admin"}))

	// ゲームを削除するエンドポイント
	router.DELETE("/delete", controllers.DeleteGameHandler, middlewares.RequireLabel([]string{"admin"}))

	// ゲームの状態を変更するエンドポイント
	router.PATCH("/start", controllers.StartGameHandler, middlewares.RequireLabel([]string{"admin"}))

	// ゲームを終了するエンドポイント
	router.PATCH("/end", controllers.EndGameHandler, middlewares.RequireLabel([]string{"admin"}))

	// チームの削除するエンドポイント
	router.DELETE("/team/delete", controllers.DeleteTeamHandler, middlewares.RequireLabel([]string{"admin"}))

	// メンバーを削除するエンドポイント
	router.DELETE("/member/delete", controllers.DeleteMemberHandler, middlewares.RequireLabel([]string{"admin"}))

	//円詳細取得
	router.GET("/circle", controllers.GetCircleDetaileHandler,middlewares.RequireAuth)

	// 円画像取得
	router.GET("/circle/image/:circle_id", controllers.GetCircleImageHandler)

	//円画像アップロード
	router.POST("/circle/image/upload", controllers.UploadCircleImageHandler, middlewares.RequireAuth)

	// 画像のリストを返すエンドポイント
	router.GET("/image/list", controllers.GetImageListHandler, middlewares.RequireAuth)

	// 終了済みゲーム一覧
	router.GET("/endgame/:user_id", controllers.GetEndGamesHandler)
	// 参加している全てのゲーム一覧取得
	router.GET("/joingame", controllers.GetJoinGamesHandler, middlewares.RequireAuth)

	// 開催中のゲーム一覧を取得する
	router.GET("/startedgame", controllers.GetStartedGamesHandler, middlewares.RequireAuth)

	// 自身が関与しているゲームの情報を取得する
	router.GET("/info/self", controllers.GetMyGamesHandler, middlewares.RequireAuth)

	return router
}
