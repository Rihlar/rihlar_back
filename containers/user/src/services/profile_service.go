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

type PrivacyInput struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Size      *int     `json:"size"`
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

	targetProfile, err := models.FindProfileById(UserID)

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

//プライバシー情報の取得
func GetPrivacyProfileService(UserID string) (*models.PrivacyProfile, error) {
	return models.FindPrivacyProfile(UserID)
}

//プライバシー情報の編集
func UpdatePrivacyProfileById(UserID string, input PrivacyInput) error {
	if UserID == "" {
		return errors.New("userID is required")
	}

	Profile, err := models.FindProfileById(UserID)
	if err != nil {
		return err
	}

	if input.Latitude != nil {
		Profile.Latitude = *input.Latitude
	}

	if input.Longitude != nil {
		Profile.Longitude = *input.Longitude
	}

	if input.Size != nil {
		Profile.Size = *input.Size
	}

	return models.UpdatePrivacyProfile(UserID, models.PrivacyProfile{
		Latitude:  Profile.Latitude,
		Longitude: Profile.Longitude,
		Size:      Profile.Size,
	})
}

// 所属地域の取得
func GetRegionProfileService(UserID string) (*models.RegionProfile, error) {
	return models.FindRegionProfile(UserID)
}

//　所属地域の編集
func UpdateRegionById(UserID string, regionID string) error {
	if UserID == "" {
		return errors.New("userID is required")
	}

	if regionID == "" {
		return errors.New("regionID is required")
	}

	profile, err := models.FindRegionProfile(UserID)
	if err != nil {
		return err
	}
	profile.RegionID = regionID

	return models.UpdateRegionProfile(UserID, profile.RegionID)
}
