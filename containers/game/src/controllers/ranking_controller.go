package controllers

import (
	"game/logger"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var rankingService = services.RnakingService{} // サービスの実体を作る。

// 自分の現在のランキング取得
func GetMyRankingHandler(c echo.Context) error {
	// ユーザーの特定する
	id := c.Param("user_id")

	// サービスに渡す
	ranking, err := rankingService.GetMyRanking(id)
	if err != nil {
		logger.PrintErr("ランキング取得エラー", ranking)
		return err
	}

	// 成功ログ
	logger.Println("Successful myRanking get.")
	// レスポンス
	c.JSON(http.StatusCreated, echo.Map{
		"helpData": ranking,
	})

	logger.Println(ranking)

	return nil
}
