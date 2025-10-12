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

	// ユーザーIDを取得する
	userId := ctx.Get("UserID").(string)

	// サービスを呼び出す
	response, err := services.ReportMovement(services.MovementArgs{
		UserID:    userId,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
	})
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	return ctx.JSON(http.StatusOK, response)
}

type AdminMovementArgs struct {
	UserID    string  `json:"user_id"`
	Steps     int64   `json:"steps`     //歩いた歩数
	Latitude  float64 `json:"latitude`  //緯度
	Longitude float64 `json:"longitude` //経度
}

// 歩いたことを報告するエンドポイント (管理者向け)
func AdminReportMovement(ctx echo.Context) error {
	// bind する
	var args []AdminMovementArgs
	if err := ctx.Bind(&args); err != nil {
		return err
	}

	// サービスを呼び出す
	for _, log := range args {
		_, err := services.ReportMovement(services.MovementArgs{
			UserID:    log.UserID,
			Steps:     log.Steps,
			Latitude:  log.Latitude,
			Longitude: log.Longitude,
		})
		if err != nil {
			logger.PrintErr(err)
			continue
			// return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to report movement"})
		}
	}

	return ctx.JSON(http.StatusOK, echo.Map{"result": "success"})
}

// 歩いたデータを記録するエンドポイント
func GetReportedMovement(ctx echo.Context) error {
	// ユーザーIDを取得する
	userId := ctx.Get("UserID").(string)

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
		"data":   movementLogs,
	})
}
