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

// 上位10位のデータを取得する
func GetRankingTop10Handler(ctx echo.Context) error {
	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

	// コンテキストからゲームID取得
	gameid := ctx.Get("GameID")
	// ID取得できてるかチェック
	if gameid == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "badRequest"})
	}

	// サービスに渡す
	ranking, err := rankingService.GetRankingTop10(userid, gameid.(string))
	if err != nil {
		logger.PrintErr("ランキング取得エラー", ranking)
		return err
	}

	// 成功ログ
	logger.Println("Successful Ranking get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": ranking,
	})
}

// ソロ用のランキングTOP10
func GetRankingTop10SoloHandler(ctx echo.Context) error {
	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

		// コンテキストからゲームID取得
	gameid := ctx.Get("GameID")
	// ID取得できてるかチェック
	if gameid == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "badRequest"})
	}

	// サービスに渡す
	ranking, err := rankingService.GetSoloRankingTop10(userid, gameid.(string))
	if err != nil {
		logger.PrintErr("ランキング取得エラー", ranking)
		return err
	}

	// 成功ログ
	logger.Println("Successful Ranking get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": ranking,
	})
}
