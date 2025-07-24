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

// 上位10位のデータを取得する
func GetRankingTop10Handler(ctx echo.Context) error {
	// ゲームのIDを取得する
	gameId := ctx.Param("game_id")

	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

	// サービスに渡す
	ranking, err := rankingService.GetRankingTop10(userid, gameId)
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
	// ゲームのIDを取得する
	gameId := ctx.Param("game_id")

	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

	// サービスに渡す
	ranking, err := rankingService.GetSoloRankingTop10(userid, gameId)
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

// 上位３位の円を取得する
func GetRankingTopHandler(ctx echo.Context) error {
	// ゲームの特定する
	id := ctx.Param("game_id")

	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

	logger.Println("UserID: ", userid)

	// サービスに渡す
	ranking, err := rankingService.GetRankingTop(userid, id)
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

// 画像取得
func GetCircleImageHandler(ctx echo.Context) error {

	// circleIDの特定　TODO: 
	id := ctx.Param("circle_id")
  	// id := ctx.Request().Header.Get("CircleID")

	// サービスに渡す
	imagePath, err := circleService.GetCircleImage(id)
	if err != nil {
		logger.PrintErr("画像パス取得エラー", imagePath)
		return err
	}

	// 成功ログ
	logger.Println("Successful imagePath get.")
	
	// レスポンス
	return ctx.File(
		imagePath,
	)
}

// 円の画像アップロード
func UploadCircleImageHandler(ctx echo.Context) error {
	// circleIDの特定　TODO:
	id := ctx.Request().Header.Get("CircleID")

	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

	logger.Println("UserID: ", userid)

	// ファイルの特定
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "画像ファイルが必要です"})
	}

	// サービスに渡す
	err = circleService.UploadImage(id, userid, fileHeader)
	if err != nil {
		logger.PrintErr("画像アップロードエラー", err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "画像アップロードエラー"})
	}

	// 成功ログ
	logger.Println("Successful image upload.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": "success",
		"circle_id": id,
	})
}

// 画像のリスト
func GetImageListHandler(ctx echo.Context) error {
	// TODO UserID の取得 (後々ミドルウェアからの取得に変更する)
	// userid := ctx.Request().Header.Get("UserID")
	userid := ctx.Get("UserID").(string)

	// サービスに渡す
	imageList, err := circleService.GetImageList(userid)
	if err != nil {
		logger.PrintErr("画像リスト取得エラー", imageList)
		return err
	}

	// 成功ログ
	logger.Println("Successful image list get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": imageList,
	})
}