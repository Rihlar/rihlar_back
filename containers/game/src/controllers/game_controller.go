package controllers

import (
	"game/logger"
	"game/models"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var gameService = services.GameService{} //サービスの実体を作る


//　管理者ゲーム作成
func CreateAdminGameHandler(ctx echo.Context) error {

	// 構造体にマッピング
	var bGame models.Game

	// バインドもしリクエストが無効ならbadRequest
	if err := ctx.Bind(&bGame); err != nil {
		logger.PrintErr("Failure to bind request.", err)

		return ctx.JSON(http.StatusBadRequest, echo.Map{
			// 空のオブジェクトをJSONとして返す
			"Data": map[string]interface{}{},
		})
	}

	err := gameService.CreateAdminGame(bGame)
		if err != nil {
		logger.PrintErr("ゲーム作成エラー", err)
		return err
	}

		// 成功ログ
	logger.Println("Successful admingame create.")
	
	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "Successful admingame create.",
	})
}