package controllers

import (
	"game/logger"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var rankingService = services.RankingService{} // サービスの実体を作る。

// 自分の現在のランキング取得
func GetMyRankingHandler(ctx echo.Context) error {
	// ユーザーの特定する
	id := ctx.Param("user_id")

	// サービスに渡す
	ranking, err := rankingService.GetMyRanking(id)
	if err != nil {
		logger.PrintErr("ランキング取得エラー", ranking)
		return err
	}

	// 成功ログ
	logger.Println("Successful myRanking get.")
	
	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": ranking,
	})
}


