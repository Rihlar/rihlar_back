package services

import (
	"errors"
	"user/models"
)

type Input struct {
	Name     string `json:"name"`           //ユーザ名
	RecordID string `json:"record_id"`      // 実績ID
	Comment  string `json:"comment"`        //コメント
	RegionID string `json:"region_id"`      // 地域ID
	SysGame  string `json:"system_game_id"` // システムゲームID
	AdmGame  string `json:"admin_game_id"`  // アドミンゲームID
}

// Profile情報を全て取得
func GetProfileService(userID string) (*models.Profile, error) {
	return models.FindProfileById(userID)
}

// Profile情報の入力部分を編集
func UpdateProfileById(UserID string, input Input) error {
	if UserID == "" {
		return errors.New("userID is required")
	}

	targetProfile,err := models.FindProfileById(UserID)

	// エラー処理
	if err != nil {
		return err
	}

	if input.Name != "" {
		targetProfile.Name = input.Name
	}

	if input.RecordID != "" {
		targetProfile.RecordID = input.RecordID
	}

	if input.Comment != "" {
		targetProfile.Comment = input.Comment
	}

	if input.RegionID != "" {
		targetProfile.RegionID = input.RegionID
	}

	if input.SysGame != "" {
		targetProfile.SysGame = input.SysGame
	}

	if input.AdmGame != "" {
		targetProfile.AdmGame = input.AdmGame
	}

	return models.UpdateProfile(UserID, *targetProfile)
}
