package controllers

import (
	"gcore/logger"
	"gcore/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MovementArgs struct {
	Steps     int64   `json:"steps`     //歩いた歩数
	Latitude  float64 `json:"latitude`  //緯度
	Longitude float64 `json:"longitude` //経度
}

// 歩いたことを報告するエンドポイント
func ReportMovement(ctx echo.Context) error {
	// bind する
	args := MovementArgs{}
	if err := ctx.Bind(&args); err != nil {
		return err
	}

	// TODO 後ほどミドルウェアからの取得に変更する
	userId := ctx.Request().Header.Get("UserID")

	// サービスを呼び出す
	if err := services.ReportMovement(services.MovementArgs{
		UserID:    userId,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
	}); err != nil {
		logger.PrintErr(err)
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"result": "success",
	})
}

// 歩いたデータを記録するエンドポイント
func GetReportedMovement(ctx echo.Context) error {
	// TODO 後ほどミドルウェアからの取得に変更する
	userId := ctx.Request().Header.Get("UserID")
	gameId := ctx.Request().Header.Get("GameID")

	// サービスを呼び出す
	movementLogs, err := services.GetReportedMovement(gameId, userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "internal server error"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"result": "success",
		"data": movementLogs,
	})
}
