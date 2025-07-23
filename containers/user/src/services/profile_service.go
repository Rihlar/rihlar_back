package services

import (
	"errors"
	"user/models"
)

// 基本プロフィールの構造体
type Input struct {
	Name     string `json:"name"`           //ユーザ名
	RecordID string `json:"record_id"`      // 実績ID
	Comment  string `json:"comment"`        //コメント
	RegionID string `json:"region_id"`      // 地域ID
	SysGame  string `json:"system_game_id"` // システムゲームID
	AdmGame  string `json:"admin_game_id"`  // アドミンゲームID
}

// プライバシー情報の構造体
type PrivacyInput struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Size      *int     `json:"size"`
}

// Profile情報を全て取得
func GetProfileService(userID string) (*models.Profile, error) {
	return models.FindProfileById(userID)
}

//Profileを作成
func CreateProfileService(userid string,input Input) (string, error){
	//構造体にinputを格納
	profile := models.Profile{
		Name: input.Name,
		RecordID: input.RecordID,
		Comment:  input.Comment,
		RegionID: input.RegionID,
		SysGame: input.SysGame,
		AdmGame: input.AdmGame,
	}

	return models.CreateProfile(userid,profile)
}

// Profile情報の入力部分を編集
func UpdateProfileById(UserID string, input Input) error {
	//ユーザーIDがからの時
	if UserID == "" {
		return errors.New("userID is required")
	}

	targetProfile, err := models.FindProfileById(UserID)

	// エラー処理
	if err != nil {
		return err
	}

	//名前が空文字でない時のみ代入
	if input.Name != "" {
		targetProfile.Name = input.Name
	}

	//実績IDが空文字でない時のみ代入
	if input.RecordID != "" {
		targetProfile.RecordID = input.RecordID
	}

	//コメントが空文字でない時のみ代入
	if input.Comment != "" {
		targetProfile.Comment = input.Comment
	}

	//地域情報が空文字でない時のみ代入
	if input.RegionID != "" {
		targetProfile.RegionID = input.RegionID
	}

	//システムゲームIDが空文字でない時のみ代入
	if input.SysGame != "" {
		targetProfile.SysGame = input.SysGame
	}

	//アドミンゲームIDが空文字でない時のみ代入
	if input.AdmGame != "" {
		targetProfile.AdmGame = input.AdmGame
	}

	return models.UpdateProfile(UserID, *targetProfile)
}

// プライバシー情報の取得
func GetPrivacyProfileService(UserID string) (*models.PrivacyProfile, error) {
	return models.FindPrivacyProfile(UserID)
}

// プライバシー情報の編集
func UpdatePrivacyProfileById(UserID string, input PrivacyInput) error {

	//UserIDが空文字ならエラー返す
	if UserID == "" {
		return errors.New("userID is required")
	}

	//プライバシー情報を取得し、失敗ならエラーを返す
	privacyProfile, err := models.FindPrivacyProfile(UserID)
	if err != nil {
		return err
	}

	//緯度が空文字でない時のみ代入
	if input.Latitude != nil {
		privacyProfile.Latitude = *input.Latitude
	}

	//経度が空文字でない時のみ代入
	if input.Longitude != nil {
		privacyProfile.Longitude = *input.Longitude
	}

	//サイズが空文字でない時のみ代入
	if input.Size != nil {
		privacyProfile.Size = *input.Size
	}

	//緯度経度サイズそれぞれを編集する
	return models.UpdatePrivacyProfile(UserID, models.PrivacyProfile{
		Latitude:  privacyProfile.Latitude,
		Longitude: privacyProfile.Longitude,
		Size:      privacyProfile.Size,
	})
}

// 所属地域の取得
func GetRegionProfileService(UserID string) (string, error) {
	return models.FindRegionProfile(UserID)
}

// 　所属地域の編集
func UpdateRegionById(UserID string, regionID string) error {
	//ユーザーIDが空ならエラーを返す
	if UserID == "" {
		return errors.New("userID is required")
	}

	//編集する地域情報がからならエラーを返す
	if regionID == "" {
		return errors.New("regionID is required")
	}

	//UserIDから地域情報取得を失敗したらエラーを返す
	_, err := models.FindRegionProfile(UserID)
	if err != nil {
		return err
	}

	return models.UpdateRegionProfile(UserID, regionID)
}
