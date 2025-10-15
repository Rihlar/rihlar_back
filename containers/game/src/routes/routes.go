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
	router.GET("/ranking/personal/:user_id", controllers.GetMyRankingHandler,middlewares.RequestLogger())

	// 上位３位の円を取得する
	router.GET("/ranking/top/:game_id", controllers.GetRankingTopHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// rankingTop10
	router.GET("/ranking/top10/:game_id", controllers.GetRankingTop10Handler, middlewares.RequireAuth,middlewares.RequestLogger())

	// ソロ用のランキングtop10
	router.GET("/ranking/solo/top10/:game_id", controllers.GetRankingTop10SoloHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// ゲームの一覧を取得するエンドポイント (モバイル用)
	router.GET("/allgames", controllers.GetAllGameListHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// ゲームを作成するエンドポイント
	router.POST("/create", controllers.CreateGameHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	// ゲームに参加するエンドポイント
	router.POST("/join", controllers.JoinGameHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// ゲームの一覧を取得するエンドポイント
	router.GET("/list", controllers.GetGameListHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	// 特定のゲームのチーム一覧を取得するエンドポイント
	router.GET("/game/:game_id/teams", controllers.GetTeamsHandler, middlewares.RequireLabel([]string{"admin"}), middlewares.RequestLogger())

	// ゲームを削除するエンドポイント
	router.DELETE("/delete", controllers.DeleteGameHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	// ゲームの状態を変更するエンドポイント
	router.PATCH("/start", controllers.StartGameHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	// ゲームを終了するエンドポイント
	router.PATCH("/end", controllers.EndGameHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	// チームの削除するエンドポイント
	router.DELETE("/team/delete", controllers.DeleteTeamHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	// メンバーを削除するエンドポイント
	router.DELETE("/member/delete", controllers.DeleteMemberHandler, middlewares.RequireLabel([]string{"admin"}),middlewares.RequestLogger())

	//円詳細取得
	router.GET("/circle", controllers.GetCircleDetaileHandler,middlewares.RequireAuth,middlewares.RequestLogger())

	// 円画像取得 (ここは認証なし)
	router.GET("/circle/image/:circle_id", controllers.GetCircleImageHandler)

	//円画像アップロード
	router.POST("/circle/image/upload", controllers.UploadCircleImageHandler, middlewares.RequireAuth)

	// 画像のリストを返すエンドポイント
	router.GET("/image/list", controllers.GetImageListHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// 終了済みゲーム一覧
	router.GET("/endgame", controllers.GetEndGamesHandler, middlewares.RequireAuth,middlewares.RequestLogger())
	// 参加している全てのゲーム一覧取得
	router.GET("/joingame", controllers.GetJoinGamesHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// 開催中のゲーム一覧を取得する
	router.GET("/startedgame", controllers.GetStartedGamesHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// 自身が関与しているゲームの情報を取得する
	router.GET("/info/self", controllers.GetMyGamesHandler, middlewares.RequireAuth,middlewares.RequestLogger())

	// アイテム詳細取得
	router.GET("/item", controllers.GetItemDeteileHandler)

	// 所持アイテム取得
	router.GET("/item/box", controllers.GetItemBoxHandler)

	// ガチャ
	router.GET("item/gacha", controllers.GetItemGachaHandler)

	// 管理者用API
	adminGroup := router.Group("/admin")
	adminGroup.Use(middlewares.RequireLabel([]string{"admin"})) // 管理者ミドルウェアを適用

	// ユーザーが参加したゲーム一覧を取得するエンドポイント (管理者用)
	adminGroup.GET("/users/:user_id/games", controllers.GetUserParticipatedGamesHandler, middlewares.RequestLogger())

	// 特定のゲームの行動ログを取得するエンドポイント (管理者用)
	adminGroup.GET("/games/:game_id/movement_logs/:user_id", controllers.GetMovementLogsHandler, middlewares.RequestLogger())

	// ゲームにユーザーを追加するエンドポイント (管理者用)
	adminGroup.POST("/member/join", controllers.AdminAddUserToGameHandler, middlewares.RequestLogger())

	return router
}
