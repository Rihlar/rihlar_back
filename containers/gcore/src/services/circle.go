package services

import (
	"gcore/location"
	"gcore/logger"
	"gcore/models"
	"gcore/utils"

	"gorm.io/gorm"
)

type CreateCircleArgs struct {
	UserID    string  `json:"userID`
	Steps     int64   `json:"steps`
	Latitude  float64 `json:"latitude`
	Longitude float64 `json:"longitude`
	Theme     string  `json:"theme`
}

func CreateCircle(args CreateCircleArgs) (CreateCirclesResponse, error) {
	// ユーザーのプロファイルを取得する
	profile, err := models.GetProfile(args.UserID)

	// エラー処理
	if err != nil {
		return CreateCirclesResponse{}, err
	}

	games := []models.Game{}

	// adminゲームを取得
	admGame, err := models.GetGame(profile.AdmGame)

	// エラー処理
	if err != nil {
		// 見つからない時
		logger.PrintErr("adminゲーム取得エラー", err)
	} else {
		// 見つかった時
		games = append(games, admGame)
	}

	// システムゲームを取得
	sysGame, err := models.GetGame(profile.SysGame)

	// エラー処理
	if err != nil {
		return CreateCirclesResponse{}, err
	}

	// システムゲームを追加
	games = append(games, sysGame)

	// 円のサイズを計算する
	circleSize := StepToSize(args.Steps)

	// 円を作成 (システムとadminゲーム)
	response, err := ProcessCreateCircle(GamesCreateCircleArgs{
		UserID:    args.UserID,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     games,
		CircleSize: circleSize,
		Theme: args.Theme,
	})

	// エラー処理
	if err != nil {
		return CreateCirclesResponse{}, err
	}

	return response, nil
}

type GamesCreateCircleArgs struct {
	UserID     string  `json:"userID`
	Steps      int64   `json:"steps`
	Latitude   float64 `json:"latitude`
	Longitude  float64 `json:"longitude`
	Games      []models.Game
	CircleSize float64 `json:"circleSize` //円のサイズ
	Theme      string  `json:"theme`		//円のテーマ
}


func StepToSize(steps int64) float64 {
	if steps < 999 {
		return 50
	}

	if steps < 2999 {
		return 100
	}

	if steps < 5999 {
		return 200
	}

	if steps < 9999 {
		return 300
	}

	// TODO 最終的にここで計算する (メートルを返す)
	return 400 //float64(steps * 10.0)
}

type CreateCirclesResponse struct {
	IsAdmin        bool   `json:"isAdmin`
	AdminCircleID string `json:"adminCircleID`
	SystemCircleID string `json:"systemCircleID`
}

// 円のデータを保存してチャンクを更新する関数 (円のIDリスト,エラー)
func ProcessCreateCircle(args GamesCreateCircleArgs) (CreateCirclesResponse, error) {
	// ID 生成
	ImageId, _ := utils.Genid()
	// circle のID を生成する
	ImageIdStr := "image-" + ImageId

	// 返すデータ
	response := CreateCirclesResponse{}

	for _, game := range args.Games {
		// 円を作成する関数
		// メンバーを取得する関数
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return CreateCirclesResponse{}, err
		}

		// ID 生成
		circleId, _ := utils.Genid()
		// circle のID を生成する
		circleIdStr := "circle-" + circleId

		circleData := &models.Circle{
			CircleID:  circleIdStr,
			GameID:    game.GameID,
			UserID:    member.UserID,
			TeamID:    member.TeamID,
			Size:      int(args.CircleSize),
			Level:     2,
			Latitude:  args.Latitude,
			Longitude: args.Longitude,
			ImageID:   ImageIdStr,
			Steps:     args.Steps,
			Theme:     args.Theme,
		}

		// 円を作成する
		if err := member.CreateCircle(circleData); err != nil {
			logger.PrintErr(err)
			return CreateCirclesResponse{}, err
		}

		// ID を追加
		// CircleIds = append(CircleIds, circleIdStr)

		if game.Type == 0 {
			// システムの時
			response.SystemCircleID = circleIdStr
		} else {
			// 管理ゲームの時
			response.AdminCircleID = circleIdStr
			response.IsAdmin = true
		}
	}

	return response, nil
}


type ProcessCircleArgs struct {
	UserID    string        `json:"userID`    // ユーザーID
	Latitude  float64       `json:"latitude`  // 緯度
	Longitude float64       `json:"longitude` // 経度
	Games     []models.Game `json:"games`     // 管理ゲームID
	Steps     int64         `json:"steps`     // 歩数
}

// 円をレベルアップさせる処理
func LevelUpCircle(args ProcessCircleArgs) error {
	// ゲームを回す
	for _, game := range args.Games {
		// メンバーを取得
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return err
		}

		// レベル2の円を取得
		levelTwoCircles, err := member.GetLevelTwoCircles()

		// エラー処理
		if err != nil {
			return err
		}

		logger.Println("level two circles:", levelTwoCircles)

		// 円を回す
		for _, circle := range levelTwoCircles {
			// 円をキャッシュ (ユーザーIDごとのキャッシュ)
			err := location.CacheCircle(location.CircleData{
				Center:   location.LatLng{
					Lat: circle.Latitude,
					Lng: circle.Longitude,
				},
				Radius:   float64(circle.Size),
				CircleID: circle.CircleID,
				UserID: member.UserID,
			})

			// エラー処理
			if err != nil {
				return err
			}
		}
	}


	// キャッシュから近い円を取得する
	NearCircles, err := location.GetNearCircle(args.UserID, location.LatLng{
		Lat: args.Latitude,
		Lng: args.Longitude,
	}, CircleMaxSize)

	// エラー処理
	if err != nil {
		return err
	}

	logger.Println("near circles:", NearCircles)

	// 縁が近い順に並んでいるので (円のサイズより距離が大きい円を見つけるまでループ)
	for _, NearCiecle := range NearCircles {
		// 円のデータを取得する
		circle, err := models.GetCircle(NearCiecle.CircleID)

		if err == gorm.ErrRecordNotFound {
			// 見つからない時
			// キャッシュから削除
			err := location.DeleteCircle(location.CircleData{
				CircleID: NearCiecle.CircleID,
				UserID:   args.UserID,
			})

			// エラー処理
			if err != nil {
				logger.PrintErr("delete circle error:", err)
				continue
			}
		}

		// エラー処理
		if err != nil {
			logger.PrintErr("near circle error:", err)
			continue
		}

		// レベルが2以外なら
		if circle.Level != 2 {
			// キャッシュから削除
			err := location.DeleteCircle(location.CircleData{
				CircleID: NearCiecle.CircleID,
				UserID:   args.UserID,
			})

			// エラー処理
			if err != nil {
				logger.PrintErr("delete circle error:", err)
				continue
			}
			continue
		}

		logger.Println("Size:", circle.Size, "Distance:", NearCiecle.Distance)

		if circle.Size < int(NearCiecle.Distance) {
			// 距離のほうが大きい (入っていない円を見つけた時)
			// ループを抜ける
			break
		}

		logger.Println("steps:", circle.Steps, "args steps:", args.Steps)
		// 距離のほうが小さい時 (ユーザーが入っている時)
		// 歩数を判定する
		if (args.Steps - circle.Steps) > LevelUpSteps {
			// 歩数を更新する
			if err := circle.ChangeLevel(3); err != nil {
				return err
			}
		}
	}

	return nil
}