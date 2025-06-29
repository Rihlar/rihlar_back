package models

import (
	"user/logger"

	"gorm.io/gorm"
)

// テーブル構造
type Profile struct {
	UserID    string  `gorm:"primaryKey;type:varchar(50)" json:"user_id"`        // ユーザID
	Name      string  `gorm:"type:varchar(100);default:''" json:"name"`          //ユーザ名
	RecordID  string  `gorm:"type:varchar(50)" json:"record_id"`                 // 実績ID
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

//地域設定情報のみの構造体
type RegionProfile struct {
	RegionID string `json:"region_id"`
}

func (Profile) TableName() string {
	return "Profile"
}

// IDからプロフィールを返す
func FindProfileById(userID string) (*Profile, error) {
	var profile Profile
	if err := dbconn.Where(&Profile{
		UserID:    userID,
	}).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// プロフィールを編集する
func UpdateProfile(userID string, data Profile) error {
	result := dbconn.
		Where("user_id = ?", userID).
		Updates(&data)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

//プライバシー情報を返す
func FindPrivacyProfile(UserID string) (*PrivacyProfile, error) {
	var result PrivacyProfile
	err := dbconn.Model(&Profile{}).
	Select("latitude", "longitude", "size").
	Where("user_id = ?", UserID).
	Scan(&result).Error
	return &result, err
}

//プライバシー情報を編集する
func UpdatePrivacyProfile(UserID string, data PrivacyProfile) error {
	return dbconn.Model(&Profile{}).
		Where("user_id = ?", UserID).
		Updates(data).Error
}

//地域情報を返す
func FindRegionProfile (UserID string) (*RegionProfile, error) {
	var result RegionProfile
	err := dbconn.Model(&Profile{}).
	Select("region_id").
	Where("user_id = ?", UserID).
	Scan(&result).Error
	return &result, err
}

func UpdateRegionProfile (UserID string, regionID string) error {
	return dbconn.Model(&Profile{}).
	Where("user_id = ?", UserID).
	Update("region_id", regionID).Error
}

// デバッグ用
func DebugProfile() {
	user_id := "e3abf90d-4bcf-4c3b-bbde-37694b1611b3"
	system_game_id := "f5b632cb-707d-f450-eece-f119534b724c"
	admin_game_id := "01e32526-1c27-c9a9-2b5b-d158b0f50c83"

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

	logger.Println("プロフィール取得成功")
}
