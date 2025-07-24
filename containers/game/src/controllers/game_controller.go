package controllers

import (
	"game/logger"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var gameService = services.GameService{} // サービスの実体を作る。

// 終了済みゲーム一覧
func GetEndGamesHandler(ctx echo.Context) error {
	// ユーザーの特定する
	id := ctx.Param("user_id")

	// サービスに渡す
	endGame, err := gameService.GetEndGames(id)
	if err != nil {
		logger.PrintErr("終了済みゲーム取得エラー", endGame)
		return err
	}

	// 成功ログ
	logger.Println("Successful endgames get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": endGame,
	})
}

// 参加中のゲーム取得
func GetJoinGamesHandler(ctx echo.Context) error {
	// ユーザーの特定する
	id := ctx.Param("user_id")

	// サービスに渡す
	joinGames, err := gameService.GetJoinGames(id)
	if err != nil {
		logger.PrintErr("参加ゲーム取得エラー", joinGames)
		return err
	}

	// 成功ログ
	logger.Println("Successful joinGames get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": joinGames,
	})
}

// ゲームに参加するエンドポイント
func JoinGameHandler(ctx echo.Context) error {
	// UserID を取得
	userId := ctx.Get("UserID").(string)

	gameId := ctx.Request().Header.Get("GameID")

	// サービスに渡す
	err := gameService.JoinGame(userId, gameId)
	if err != nil {
		logger.PrintErr("ゲームに参加エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful game join.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

// ゲームを作成する関数
func CreateGameHandler(ctx echo.Context) error {
	// bodyを取得
	var args services.CreateGameArgs
	if err := ctx.Bind(&args); err != nil {
		return err
	}

	// サービスに渡す
	err := gameService.CreateGame(args)
	if err != nil {
		logger.PrintErr("ゲーム作成エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful game create.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

// ゲーム一覧
func GetGameListHandler(ctx echo.Context) error {
	// サービスに渡す
	games, err := gameService.GetGameList()
	if err != nil {
		logger.PrintErr("ゲーム一覧取得エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful game list get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": games,
	})
}

// ゲームを削除する
func DeleteGameHandler(ctx echo.Context) error {
	// ゲームIDを取得
	gameId := ctx.Request().Header.Get("GameID")

	// サービスに渡す
	err := gameService.DeleteGame(gameId)
	if err != nil {
		logger.PrintErr("ゲーム削除エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful game delete.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

type ToggleGameStatusArgs struct {
	GameId string `json:"game_id"`
}

// ゲームを開始するエンドポイント
func StartGameHandler(ctx echo.Context) error {
	// ゲームIDを取得
	var args ToggleGameStatusArgs
	if err := ctx.Bind(&args); err != nil {
		return err
	}

	gameId := args.GameId

	// サービスに渡す
	err := gameService.StartGame(gameId)
	if err != nil {
		logger.PrintErr("ゲーム開始エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful game start.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

// ゲームを終了するエンドポイント
func EndGameHandler(ctx echo.Context) error {
	// ゲームIDを取得
	var args ToggleGameStatusArgs
	if err := ctx.Bind(&args); err != nil {
		return err
	}

	gameId := args.GameId

	// サービスに渡す
	err := gameService.EndGame(gameId)
	if err != nil {
		logger.PrintErr("ゲーム終了エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful game end.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

// チームを削除するエンドポイント
func DeleteTeamHandler(ctx echo.Context) error {
	// TODO 後ほどミドルウェアからの取得に変更す
	teamId := ctx.Request().Header.Get("TeamID")
	gameId := ctx.Request().Header.Get("GameID")

	// サービスに渡す
	err := gameService.DeleteTeam(gameId, teamId)
	if err != nil {
		logger.PrintErr("チーム削除エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful team delete.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

// メンバーを削除するエンドポイと
func DeleteMemberHandler(ctx echo.Context) error {
	// TODO 後ほどミドルウェアからの取得に変更す
	userId := ctx.Get("UserID").(string)
	// userId := ctx.Request().Header.Get("UserID")
	gameId := ctx.Request().Header.Get("GameID")

	// サービスに渡す
	err := gameService.DeleteMember(gameId, userId)
	if err != nil {
		logger.PrintErr("メンバー削除エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful member delete.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
	})
}

// 開催中のゲームを取得する
func GetStartedGamesHandler(ctx echo.Context) error {
	// ユーザーIDを取得する
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	games, err := gameService.GetStartedGames(userId)
	if err.Err != nil {
		logger.PrintErr("開催中のゲーム取得エラー", err.LogMessage)
		return ctx.JSON(err.Code, echo.Map{
			"Message": err.ErrMessage,
		})
	}

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": games,
	})
}

// 自身が関与しているゲームの一覧を取得する
func GetMyGamesHandler(ctx echo.Context) error {
	// ユーザーIDを取得する
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	data, err := gameService.GetMyGames(userId)
	if err != nil {
		logger.PrintErr("自身が関与しているゲーム取得エラー", err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"Message": "自身が関与しているゲーム取得エラー",
		})
	}

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": data,
	})
}
