package services

import (
	"gcore/logger"
	"gcore/models"
	"gcore/utils"
)

type CreateCircleArgs struct {
	UserID    string  `json:"userID`
	Steps     int64   `json:"steps`
	Latitude  float64 `json:"latitude`
	Longitude float64 `json:"longitude`
}

func CreateCircle(args CreateCircleArgs) (string,error) {
	// // ユーザーのプロファイルを取得する
	// profile, err := models.GetProfile(args.UserID)

	// // エラー処理
	// if err != nil {
	// 	return "", err
	// }

	// // adminゲームを取得
	// admGame, err := models.GetGame(profile.AdmGame)

	// // エラー処理
	// if err != nil {
	// 	return "", err
	// }

	// // システムゲームを取得
	// sysGame, err := models.GetGame(profile.SysGame)

	// // エラー処理
	// if err != nil {
	// 	return "", err
	// }
	return  "", nil
}

type GamesCreateCircleArgs struct {
	UserID    string  `json:"userID`
	Steps     int64   `json:"steps`
	Latitude  float64 `json:"latitude`
	Longitude float64 `json:"longitude`
	Games     []models.Game
}

func PorcessCreateCircle(args GamesCreateCircleArgs) error {
	for _, game := range args.Games {
		// 円を作成する関数
		// メンバーを取得する関数
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return err
		}

		// ID 生成
		circleId,_ := utils.Genid()
		// circle のID を生成する
		circleIdStr := "circle-" + circleId

		// ID 生成
		ImageId,_ := utils.Genid()
		// circle のID を生成する
		ImageIdStr := "image-" + ImageId

		// 円を作成する
		if err := member.CreateCircle(&models.Circle{
			CircleID:  circleIdStr,
			GameID:    game.GameID,
			UserID:    member.UserID,
			Size:      0,
			Level:     2,
			Latitude:  args.Latitude,
			Longitude: args.Longitude,
			ImageID:   ImageIdStr,
		}); err != nil {
			logger.PrintErr(err)
			return err
		}
	}

	return nil
}