package models

import (
	"errors"
	"user/logger"

	"gorm.io/gorm"
)

// テーブル構造
type Profile struct {
	UserID           string  `gorm:"primaryKey;type:varchar(50)" json:"user_id"`            // ユーザID
	Name             string  `gorm:"type:varchar(100);default:''" json:"name"`              //ユーザ名
	DisplayAchiveID1 string  `gorm:"type:varchar(50);default:''" json:"display_achive_id1"` // 実績ID1
	DisplayAchiveID2 string  `gorm:"type:varchar(50);default:''" json:"display_achive_id2"` // 実績ID2
	DisplayAchiveID3 string  `gorm:"type:varchar(50);default:''" json:"display_achive_id3"` // 実績ID3
	Comment          string  `gorm:"type:varchar(255);default:''" json:"comment"`           // ユーザコメント（デフォルト空白）
	Latitude         float64 `gorm:"type:double;default:0" json:"latitude"`                 // 緯度（デフォルト0）
	Longitude        float64 `gorm:"type:double;default:0" json:"longitude"`                // 経度（デフォルト0）
	Size             int     `gorm:"default:0" json:"size"`                                 // サイズ（デフォルト0）
	RegionID         string  `gorm:"type:varchar(50);default:''" json:"region_id"`          // 地域ID
	SysGame          string  `gorm:"type:varchar(50);default:''" json:"system_game_id"`     // システムゲームID
	AdmGame          string  `gorm:"type:varchar(50);default:''" json:"admin_game_id"`      // アドミンゲームID
}

// プライバシー情報のみの構造体
type PrivacyProfile struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Size      int     `json:"size"`
}

// 画面に設定する実績情報のみの構造体
type AchiveProfile struct {
	DisplayAchiveID1 string `json:"display_achive_id1"` // 実績ID1
	DisplayAchiveID2 string `json:"display_achive_id2"` // 実績ID2
	DisplayAchiveID3 string `json:"display_achive_id3"` // 実績ID3
}

func (Profile) TableName() string {
	return "Profile"
}

// デバッグ用
func DebugProfile() {
	user_id := "userid-e3abf90d-4bcf-4c3b-bbde-37694b1611b3"
	system_game_id := "systemgame-f5b632cb-707d-f450-eece-f119534b724c"
	admin_game_id := "admingame-01e32526-1c27-c9a9-2b5b-d158b0f50c83"
	display_achive_ID1 := "achive-d3c9d4be-84c2-5b4d-b116-088a8c8680ab"
	display_achive_ID2 := "achive-86548e0a-41ab-9617-c248-094e2de7f4f4"
	display_achive_ID3 := "achive-46e37298-1052-a002-4cc4-1567df59bae4"

	//書き込み
	result := dbconn.Save(&Profile{
		UserID:           user_id,
		Name:             "山田太郎",
		DisplayAchiveID1: display_achive_ID1,
		DisplayAchiveID2: display_achive_ID2,
		DisplayAchiveID3: display_achive_ID3,
		Comment:          "よろしくお願いします。",
		Latitude:         35.23,
		Longitude:        135.25,
		Size:             100,
		RegionID:         "関東地方",
		SysGame:          system_game_id,
		AdmGame:          admin_game_id,
	})

	//エラー処理
	if result.Error != nil {
		logger.PrintErr("プロフィール保存エラー", result.Error)
		return
	}

	logger.PrintErr("プロフィール保存成功")

	//取得コード
	returnData := Profile{}

	//取得する
	result = dbconn.Where(&Profile{
		UserID: user_id,
	}).First(&returnData)

	if result.Error != nil {
		logger.PrintErr("プロフィール取得エラー", result.Error)
		return
	}

	logger.PrintErr("プロフィール取得成功")
}

/*
	プロフィール全体関連
*/

// IDからプロフィールを返す
func FindProfileById(userID string) (*Profile, error) {
	var profile Profile

	//userIDからのprofile一件検索
	result := dbconn.Where("user_id = ?", userID).First(&profile)

	//エラー
	if result.Error != nil {
		logger.PrintErr("プロフィール取得エラー", result.Error)
		return nil, result.Error
	}

	//profile全項目を返す
	return &profile, nil
}

// プロフィールの作成
func CreateProfile(userid string, data Profile) (string, error) {
	//uuid格納用に整形
	data.UserID = userid

	//ユーザー作成
	result := dbconn.Save(&data)

	//作成エラー
	if result.Error != nil {
		logger.PrintErr("ユーザー新規作成エラー", result.Error)
		return "", result.Error
	}

	//userIDを返す
	return data.UserID, nil
}

// プロフィールを編集する
func UpdateProfile(userID string, data Profile) error {

	//userIDからprofile編集
	result := dbconn.Model(&Profile{}).Where("user_id = ?", userID).Updates(data)

	//失敗したらエラーを返す
	if result.Error != nil {
		logger.PrintErr("プロフィール編集エラー", result.Error)
		return result.Error
	}

	//0件ならNotFoundを返す
	if result.RowsAffected == 0 {
		logger.PrintErr("変更レコード0件エラー")
		return errors.New("no rows updated")
	}

	return nil
}

// プロファイルが存在してるか判定する
func ExistProfile(userID string) (bool, error) {
	var profile Profile

	//userIDからprofile一件検索
	result := dbconn.Where("user_id = ?", userID).First(&profile)

	// gorm record not found
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	//失敗したらエラーを返す
	if result.Error != nil {
		logger.PrintErr("プロフィール取得エラー", result.Error)
		return false, result.Error
	}

	//0件ならfalseを返す
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

/*
	実績関連
*/

func FindAchiveProfile(userID string) (*AchiveProfile, error) {

	profile, err := FindProfileById(userID)

	if err != nil {
		logger.PrintErr("実績取得エラー", err)
		return nil, err
	}

	return &AchiveProfile{
		DisplayAchiveID1: profile.DisplayAchiveID1,
		DisplayAchiveID2: profile.DisplayAchiveID2,
		DisplayAchiveID3: profile.DisplayAchiveID3,
	}, nil
}

// 実績情報を編集する
func UpdateAchiveProfile(userID string, data AchiveProfile) error {
	result := dbconn.Model(&Profile{}).Where("user_id = ?", userID).Updates(&data)

	//エラーを返す
	if result.Error != nil {
		logger.PrintErr("実績編集エラー", result.Error)
		return result.Error
	}

	//0件エラー
	if result.RowsAffected == 0 {
		logger.PrintErr("変更レコード0件エラー")
		return errors.New("no rows updated")
	}

	return nil
}

/*
	プライバシー関連
*/

// プライバシー情報を返す
func FindPrivacyProfile(userID string) (*PrivacyProfile, error) {
	//ProfileをIDから全項目一件検索
	profile, err := FindProfileById(userID)

	//エラーならnilを返す
	if err != nil {
		logger.PrintErr("プライバシー取得エラー", err)
		return nil, err
	}

	//成功したら、緯度,経度,サイズを返す
	return &PrivacyProfile{
		Latitude:  profile.Latitude,
		Longitude: profile.Longitude,
		Size:      profile.Size,
	}, nil
}

// プライバシー情報を編集する
func UpdatePrivacyProfile(userID string, data PrivacyProfile) error {
	result := dbconn.Model(&Profile{}).Where("user_id = ?", userID).Updates(data)

	//エラーを返す
	if result.Error != nil {
		logger.PrintErr("プライバシー編集エラー", result.Error)
		return result.Error
	}

	//0件エラー
	if result.RowsAffected == 0 {
		logger.PrintErr("変更レコード0件エラー")
		return errors.New("no rows updated")
	}

	return nil
}

/*
	地域情報関連
*/

// 地域情報を返す
func FindRegionProfile(userID string) (string, error) {
	//ProfileをIDから全項目一件検索
	profile, err := FindProfileById(userID)

	//エラーなら空文字を返す
	if err != nil {
		logger.PrintErr("地域情報取得エラー", err)
		return "", err
	}

	//Regionを返す
	return profile.RegionID, nil
}

// 地域情報を編集する
func UpdateRegionProfile(userID string, regionID string) error {
	result := dbconn.Model(&Profile{}).Where("user_id = ?", userID).Update("region_id", regionID)

	//編集時エラー
	if result.Error != nil {
		logger.PrintErr("地域情報編集エラー", result.Error)
		return result.Error
	}

	//編集0件エラー
	if result.RowsAffected == 0 {
		logger.PrintErr("変更レコード0件エラー")
		return errors.New("no rows updated")
	}

	return nil
}
