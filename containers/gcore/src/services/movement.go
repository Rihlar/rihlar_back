package services

import (
	"gcore/logger"
	"gcore/models"
)

// 歩いたことを報告する引数
type MovementArgs struct {
	UserID    string  `json:"userID`	// ユーザーID
	Steps     int64   `json:"steps`		// 歩数
	Latitude  float64 `json:"latitude`	// 緯度
	Longitude float64 `json:"longitude`	// 経度
}

type SaveMovementLogArgs struct {
	UserID    string  `json:"userID`	// ユーザーID
	SystemID  string  `json:"systemID`	// システムゲームID
	AdminID   string  `json:"adminID`	// 管理ゲームID
	Steps     int64   `json:"steps`		// 歩数
	Latitude  float64 `json:"latitude`	// 緯度
	Longitude float64 `json:"longitude`	// 経度
}

func ReportMovement(args MovementArgs) error {
	// プロファイルを取得する
	profile,err := models.GetProfile(args.UserID)

	// エラー処理
	if err != nil {
		return err
	}

	// システムゲームを取得
	sysGame,err := models.GetGame(profile.SysGame)

	// エラー処理
	if err != nil {
		return err
	}

	// adminゲームを取得
	admGame,err := models.GetGame(profile.AdmGame)

	// エラー処理
	if err != nil {
		return err
	}

	// 行動記録を保存する
	err = SaveMovementLog(SaveMovementLogArgs{
		UserID:    args.UserID,
		SystemID:  sysGame.GameID,
		AdminID:   admGame.GameID,
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
	})

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

// 歩いたログを記録する関数
func SaveMovementLog(args SaveMovementLogArgs) error {
	logger.Println("SaveMovementLog", args)

	// admin のメンバーを取得
	adminMember, err := models.GetMemberByUserID(args.AdminID, args.UserID)

	// エラー処理
	if err != nil {
		return err
	}

	// system のメンバーを取得
	systemMember, err := models.GetMemberByUserID(args.SystemID, args.UserID)

	// エラー処理
	if err != nil {
		return err
	}

	// admin を保存
	err = adminMember.SaveMovementLog(args.Latitude, args.Longitude, args.Steps)

	// エラー処理
	if err != nil {
		return err
	}

	// system を保存
	err = systemMember.SaveMovementLog(args.Latitude, args.Longitude, args.Steps)

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}