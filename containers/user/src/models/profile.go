package models

import (
	"errors"
	"user/logger"

	"gorm.io/gorm"
)

// テーブル構造
type Profile struct {
	UserID    string  `gorm:"primaryKey;type:varchar(50)" json:"user_id"`        // ユーザID
	Name      string  `gorm:"type:varchar(100);default:''" json:"name"`          //ユーザ名
	RecordID  string  `gorm:"type:varchar(50);default''" json:"record_id"`                 // 実績ID
	Comment   string  `gorm:"type:varchar(255);default:''" json:"comment"`       // ユーザコメント（デフォルト空白）
	Latitude  float64 `gorm:"type:double;default:0" json:"latitude"`             // 緯度（デフォルト0）
	Longitude float64 `gorm:"type:double;default:0" json:"longitude"`            // 経度（デフォルト0）
	Size      int     `gorm:"default:0" json:"size"`                             // サイズ（デフォルト0）
	RegionID  string  `gorm:"type:varchar(50);default:''" json:"region_id"`      // 地域ID
	SysGame   string  `gorm:"type:varchar(50);default:''" json:"system_game_id"` // システムゲームID
	AdmGame   string  `gorm:"type:varchar(50);default:''" json:"admin_game_id"`  // アドミンゲームID
}

//プライバシー情報のみの構造体
type PrivacyProfile struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Size      int     `json:"size"`
}

func (Profile) TableName() string {
	return "Profile"
}

// デバッグ用
func DebugProfile() {
	user_id := "userid-e3abf90d-4bcf-4c3b-bbde-37694b1611b3"
	system_game_id := "systemgame-f5b632cb-707d-f450-eece-f119534b724c"
	admin_game_id := "admingame-01e32526-1c27-c9a9-2b5b-d158b0f50c83"

	//書き込み
	result := dbconn.Save(&Profile{
		UserID:    user_id,
		Name:      "山田太郎",
		RecordID:  "第一回優勝",
		Comment:   "よろしくお願いします。",
		Latitude:  35.23,
		Longitude: 135.25,
		Size:      100,
		RegionID:  "関東地方",
		SysGame:   system_game_id,
		AdmGame:   admin_game_id,
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

// IDからプロフィールを返す
func FindProfileById(userID string) (*Profile, error) {
	var profile Profile

	//userIDからのprofile一件検索
	result := dbconn.Where("user_id = ?", userID).First(&profile)

	//エラー
	if result.Error != nil{
		logger.PrintErr("プロフィール取得エラー",result.Error)
		return nil, result.Error
	}

	//profile全項目を返す
	return &profile, nil
}

//プロフィールの作成
func CreateProfile(userid string,data Profile) (string, error) {
	//uuid格納用に整形
	data.UserID = userid
	
	//ユーザー作成
	result := dbconn.Save(&data)

	//作成エラー
	if result.Error != nil{
		logger.PrintErr("ユーザー新規作成エラー",result.Error)
		return "", result.Error
	}

	//userIDを返す
	return data.UserID, nil
}

// プロフィールを編集する
func UpdateProfile(userID string, data Profile) error {

	//userIDからprofile編集
	result := dbconn.Where("user_id", userID).Save(&data)

	//失敗したらエラーを返す
	if result.Error != nil {
		logger.PrintErr("プロフィール編集エラー",result.Error)
		return result.Error
	}

	//0件ならNotFoundを返す
	if result.RowsAffected == 0 {
		logger.PrintErr("変更レコード0件エラー")
		return gorm.ErrRecordNotFound
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
		logger.PrintErr("プロフィール取得エラー",result.Error)
		return false, result.Error
	}

	//0件ならfalseを返す
	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

//プライバシー情報を返す
func FindPrivacyProfile(userID string) (*PrivacyProfile, error) {
	//ProfileをIDから全項目一件検索
	Profile, err := FindProfileById(userID)

	//エラーならnilを返す
	if err != nil{
		logger.PrintErr("プライバシー取得エラー",err)
		return nil, err
	}

	//成功したら、緯度,経度,サイズを返す
	return &PrivacyProfile{
		Latitude: Profile.Latitude,
		Longitude: Profile.Longitude,
		Size: Profile.Size,
	}, nil
}

//プライバシー情報を編集する
func UpdatePrivacyProfile(userID string, data PrivacyProfile) error {
	result := dbconn.Model(&Profile{}).Where("user_id = ?", userID).Updates(data)

	//エラーを返す
	if result.Error != nil {
		logger.PrintErr("プライバシー編集エラー",result.Error)
		return result.Error
	}

	//0件エラー
	if result.RowsAffected == 0{
		logger.PrintErr("変更レコード0件エラー")
		return gorm.ErrRecordNotFound
	}

	return nil
}

//地域情報を返す
func FindRegionProfile (userID string) (string, error) {
	//ProfileをIDから全項目一件検索
	profile, err := FindProfileById(userID)

	//エラーなら空文字を返す
	if err != nil {
		logger.PrintErr("地域情報取得エラー",err)
		return "", err
	}
	
	//Regionを返す
	return profile.RegionID, nil
}

//地域情報を編集する
func UpdateRegionProfile (userID string, regionID string) error {
	result := dbconn.Model(&Profile{}).Where("user_id = ?", userID).Update("region_id", regionID)

	//編集時エラー
	if result.Error != nil {
		logger.PrintErr("地域情報編集エラー",result.Error)
		return result.Error
	}

	//編集0件エラー
	if result.RowsAffected == 0 {
		logger.PrintErr("変更レコード0件エラー")
		return gorm.ErrRecordNotFound
	}

	return nil
}

