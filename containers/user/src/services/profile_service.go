package services

import (
	"errors"
	"time"
	"user/models"
	"user/utils"
)

// 基本プロフィールの構造体
type Input struct {
	Name     string `json:"name"`           //ユーザ名
	Comment  string `json:"comment"`        //コメント
	RegionID string `json:"region_id"`      // 地域ID
	SysGame  string `json:"system_game_id"` // システムゲームID
	AdmGame  string `json:"admin_game_id"`  // アドミンゲームID
}

// 　ユーザの実績の構造体
type AchiveInput struct {
	DisplayAchiveID1 *string `json:"display_achive_id1"` // 実績ID1
	DisplayAchiveID2 *string `json:"display_achive_id2"` // 実績ID2
	DisplayAchiveID3 *string `json:"display_achive_id3"` // 実績ID3
}

// プライバシー情報の構造体
type PrivacyInput struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Size      *int     `json:"size"`
}

/*
	プロフィール全体関連
*/

// Profile情報を全て取得
func GetProfileService(userID string) (*models.Profile, error) {
	return models.FindProfileById(userID)
}

// Profileを作成
func CreateProfileService(userid string, input Input) (string, error) {
	// プロファイルが存在するかチェック
	isExist, err := models.ExistProfile(userid)

	// エラー処理
	if err != nil {
		return "", err
	}

	if isExist {
		// 存在する時
		return "", errors.New("profile already exists")
	}

	// ゲームのIDを生成
	gameId, _ := utils.Genid()

	// ゲーム作成
	err = models.CreateGame(models.Game{
		GameID:    "sysgame-" + gameId,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Flag:      0,
		Type:      0,
		Status:    1,
		RegionID:  "",
	})

	// エラー処理
	if err != nil {
		return "", err
	}

	// チームIDを生成
	teamId, _ := utils.Genid()

	// メンバーを追加する
	err = models.DebugAddMember("sysgame-"+gameId, "teamid-"+teamId, userid)

	// エラー処理
	if err != nil {
		return "", err
	}

	// 存在しない時
	//構造体にinputを格納
	profile := models.Profile{
		Name:     input.Name,
		Comment:  input.Comment,
		RegionID: input.RegionID,
		SysGame:  "sysgame-" + gameId,
		AdmGame:  "",
	}

	// エラー処理
	_, err = models.CreateProfile(userid, profile)

	// エラー処理
	if err != nil {
		return "", err
	}

	return userid, nil
}

// Profile情報の入力部分を編集
func UpdateProfileById(userID string, input Input) error {
	//ユーザーIDがからの時
	if userID == "" {
		return errors.New("userID is required")
	}

	targetProfile, err := models.FindProfileById(userID)

	// エラー処理
	if err != nil {
		return err
	}

	//名前が空文字でない時のみ代入
	if input.Name != "" {
		targetProfile.Name = input.Name
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

	return models.UpdateProfile(userID, *targetProfile)
}

/*
	実績関連
*/

func GetAchiveProfile(userID string) (*models.AchiveProfile, error) {
	return models.FindAchiveProfile(userID)
}

func UpdateAchiveProfile(userID string, input AchiveInput) error {
	//UserIDが空文字ならエラー返す
	if userID == "" {
		return errors.New("userID is required")
	}

	//存在確認
	AchiveProfile, err := models.FindAchiveProfile(userID)

	//エラー処理
	if err != nil {
		return err
	}

	// 3つの実績が空文字でない時のみ代入
	if input.DisplayAchiveID1 != nil {
		AchiveProfile.DisplayAchiveID1 = *input.DisplayAchiveID1
	}

	if input.DisplayAchiveID2 != nil {
		AchiveProfile.DisplayAchiveID2 = *input.DisplayAchiveID2
	}

	if input.DisplayAchiveID3 != nil {
		AchiveProfile.DisplayAchiveID3 = *input.DisplayAchiveID3
	}

	//　実行結果を返す
	return models.UpdateAchiveProfile(userID, models.AchiveProfile{
		DisplayAchiveID1: AchiveProfile.DisplayAchiveID1,
		DisplayAchiveID2: AchiveProfile.DisplayAchiveID2,
		DisplayAchiveID3: AchiveProfile.DisplayAchiveID3,
	})
}

/*
	プライバシー関連
*/

// プライバシー情報の取得
func GetPrivacyProfileService(userID string) (*models.PrivacyProfile, error) {
	return models.FindPrivacyProfile(userID)
}

// プライバシー情報の編集
func UpdatePrivacyProfileById(userID string, input PrivacyInput) error {

	//UserIDが空文字ならエラー返す
	if userID == "" {
		return errors.New("userID is required")
	}

	//プライバシー情報を取得し、失敗ならエラーを返す
	privacyProfile, err := models.FindPrivacyProfile(userID)
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
	return models.UpdatePrivacyProfile(userID, models.PrivacyProfile{
		Latitude:  privacyProfile.Latitude,
		Longitude: privacyProfile.Longitude,
		Size:      privacyProfile.Size,
	})
}

/*
	地域情報関連
*/

// 所属地域の取得
func GetRegionProfileService(userID string) (string, error) {
	return models.FindRegionProfile(userID)
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

// 全てのProfile情報を取得
func GetAllProfilesService() ([]models.Profile, error) {
	return models.GetAllProfiles()
}
