package controllers

import (
	"game/logger"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var circleService = services.CircleService{} //サービスの実体を作る

// 円の詳細取得
func GetCircleDeteileHandler(ctx echo.Context) error {
	// サークルIDの特定する
	id := ctx.Param("circle_id")

	// 円の詳細取得
	circle, err := circleService.GetCircleDeteile(id)
		if err != nil {
		logger.PrintErr("円取得エラー", circle)
		return err
	}

	// 成功ログ
	logger.Println("Successful circleDeteil get.")
	
	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": circle,
	})

}