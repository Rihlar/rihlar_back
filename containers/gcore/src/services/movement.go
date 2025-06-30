package services

import (
	"gcore/logger"
	"gcore/models"
)

// 歩いたことを報告する引数
type MovementArgs struct {
	UserID    string  `json:"userID`    // ユーザーID
	Steps     int64   `json:"steps`     // 歩数
	Latitude  float64 `json:"latitude`  // 緯度
	Longitude float64 `json:"longitude` // 経度
}

type SaveMovementLogArgs struct {
	UserID    string        `json:"userID`    // ユーザーID
	SystemID  string        `json:"systemID`  // システムゲームID
	Games     []models.Game `json:"games`     // 管理ゲームID
	Steps     int64         `json:"steps`     // 歩数
	Latitude  float64       `json:"latitude`  // 緯度
	Longitude float64       `json:"longitude` // 経度
}

type ProcessChunkArgs struct {
	UserID    string        `json:"userID`    // ユーザーID
	Latitude  float64       `json:"latitude`  // 緯度
	Longitude float64       `json:"longitude` // 経度
	Games     []models.Game `json:"games`     // 管理ゲームID
}


func ReportMovement(args MovementArgs) error {
	// プロファイルを取得する
	profile, err := models.GetProfile(args.UserID)

	// エラー処理
	if err != nil {
		return err
	}

	// システムゲームを取得
	sysGame, err := models.GetGame(profile.SysGame)

	// エラー処理
	if err != nil {
		return err
	}

	// adminゲームを取得
	admGame, err := models.GetGame(profile.AdmGame)

	// エラー処理
	if err != nil {
		return err
	}

	// 行動記録を保存する
	err = SaveMovementLog(SaveMovementLogArgs{
		UserID:    args.UserID,
		SystemID:  sysGame.GameID,
		Games:     []models.Game{admGame, sysGame},
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
	})

	// エラー処理
	if err != nil {
		return err
	}

	// チャンクに対しての処理
	err = ProcessChunk(ProcessChunkArgs{
		UserID:    args.UserID,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     []models.Game{admGame, sysGame},
	})

	// エラー処理
	if err != nil {
		return err
	}

	// 円のレベルアップ処理
	err = LevelUpCircle(ProcessCircleArgs{
		UserID:    args.UserID,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     []models.Game{admGame, sysGame},
	})

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}


// 歩いたログを記録する関数
func SaveMovementLog(args SaveMovementLogArgs) error {
	// ゲームを回す
	for _, game := range args.Games {
		// メンバーオブジェクト取得
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return err
		}

		// 歩いたログを保存する
		if err := member.SaveMovementLog(args.Latitude, args.Longitude, args.Steps); err != nil {
			logger.PrintErr(err)
			return err
		}
	}

	return nil
}

// チャンクに対しての処理 (Level1 などを実行する)
func ProcessChunk(args ProcessChunkArgs) error {
	// ゲームを回す
	for _, game := range args.Games {
		// 一番近いチャンクを取得
		chunk,err := game.GetChunkByLatLon(args.Latitude, args.Longitude)

		// エラー処理
		if err != nil {
			return err
		}

		logger.Println("near chunk:", chunk)

		// チャンクのレベルが0か1なら
		if chunk.Level == 0 || chunk.Level == 1 {
			// チャンクを更新する
			if err := chunk.ChangeLevel(1); err != nil {
				logger.PrintErr(err)
				return err
			}

			// オーナーも変更する
			if err := chunk.ChangeOwner(args.UserID); err != nil {
				logger.PrintErr(err)
				return err
			}
		}
	}

	return nil
}