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

func CreateCircle(args CreateCircleArgs) ([]string,error) {
	// ユーザーのプロファイルを取得する
	profile, err := models.GetProfile(args.UserID)

	// エラー処理
	if err != nil {
		return []string{}, err
	}

	// adminゲームを取得
	admGame, err := models.GetGame(profile.AdmGame)

	// エラー処理
	if err != nil {
		return []string{}, err
	}

	// システムゲームを取得
	sysGame, err := models.GetGame(profile.SysGame)

	// エラー処理
	if err != nil {
		return []string{}, err
	}

	// 円を作成
	circleIds, err := WriteCiecleData(GamesCreateCircleArgs{
		UserID:    args.UserID,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     []models.Game{admGame, sysGame},
	})

	// エラー処理
	if err != nil {
		return []string{}, err
	}

	// 円を更新
	if err := ProcessCircleChunk(GamesCreateCircleArgs{
		UserID:    args.UserID,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     []models.Game{admGame, sysGame},
	}); err != nil {
		return []string{}, err
	}

	return circleIds, nil
}


type GamesCreateCircleArgs struct {
	UserID    string  `json:"userID`
	Steps     int64   `json:"steps`
	Latitude  float64 `json:"latitude`
	Longitude float64 `json:"longitude`
	Games     []models.Game
}

func ProcessCircleChunk(args GamesCreateCircleArgs) error {
	// ゲームを回す
	for _, game := range args.Games {
		// メンバーを取得
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return err
		}

		// 円形にチャンクを取得する
		chunks, err := game.GetCircleChunkByLatLon(args.Latitude, args.Longitude, StepToSize(args.Steps))

		// エラー処理
		if err != nil {
			return err
		}

		for _, chunk := range chunks {
			// チャンクのレベルが2 以下なら
			if chunk.Level <= 2 {
				// チャンクを更新する
				if err := chunk.ChangeLevel(2); err != nil {
					logger.PrintErr(err)
					return err
				}

				// オーナーも変更する
				if err := chunk.ChangeOwner(member.UserID); err != nil {
					logger.PrintErr(err)
					return err
				}
			}
		}
	}

	return nil
}

func StepToSize(steps int64) float64 {
	return float64(500) //float64(steps * 10.0)
}

// 円のデータを保存する関数 (円のIDリスト,エラー)
func WriteCiecleData(args GamesCreateCircleArgs) ([]string,error) {
	imageIds := []string{}

	for _, game := range args.Games {
		// 円を作成する関数
		// メンバーを取得する関数
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return []string{}, err
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
			return []string{}, err
		}

		// ID を追加
		imageIds = append(imageIds, ImageIdStr)
	}

	return imageIds, nil
}
