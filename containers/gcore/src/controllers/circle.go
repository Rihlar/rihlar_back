package controllers

import (
	"gcore/logger"
	"gcore/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// import (
// 	"gcore/services"

// 	"github.com/labstack/echo/v4"
// )


type CircleArgs struct {
	Latitude  float64 `json:"latitude`  //緯度
	Longitude float64 `json:"longitude` //経度
	Steps     int64   `json:"steps`     //歩いた歩数
	Theme 	  string  `json:"theme`		//テーマ
}

// 円を作成する関数
func CreateCircle(ctx echo.Context) error {
	// bind する
	args := CircleArgs{}
	if err := ctx.Bind(&args); err != nil {
		return err
	}

	// ユーザーIDを取得する
	userId := ctx.Get("UserID").(string)
	
	// TODO デバッグ用にヘッダからユーザーIDを取得する
	// userId := ctx.Request().Header.Get("UserID")

	// 円を作成する
	circleIds, err := services.CreateCircle(services.CreateCircleArgs{
		UserID:    userId,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Theme:     args.Theme,
	})

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"result": "success",
		"circleIds": circleIds,
	})
}