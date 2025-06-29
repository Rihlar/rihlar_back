package controllers

// import (
// 	"gcore/services"

// 	"github.com/labstack/echo/v4"
// )


type CircleArgs struct {
	Latitude  float64 `json:"latitude`  //緯度
	Longitude float64 `json:"longitude` //経度
	Steps     int64   `json:"steps`     //歩いた歩数
}

// // 円を作成する関数
// func CreateCircle(ctx echo.Context) error {
// 	// bind する
// 	args := CircleArgs{}
// 	if err := ctx.Bind(&args); err != nil {
// 		return err
// 	}

// 	// TODO 後ほどミドルウェアからの取得に変更する
// 	userId := ctx.Request().Header.Get("UserID")

// 	// 円を作成する
// 	services.CreateCircle(services.CircleArgs{
// 		UserID:    userId,
// 		Steps:     args.Steps,
// 		Latitude:  args.Latitude,
// 		Longitude: args.Longitude,
// 	})

// 	return ctx.JSON(http.StatusOK, echo.Map{
// 		"result": "success",
// 	})
// }