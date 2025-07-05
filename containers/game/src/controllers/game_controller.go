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
	// TODO 後ほどミドルウェアからの取得に変更する
	userId := ctx.Request().Header.Get("UserID")
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