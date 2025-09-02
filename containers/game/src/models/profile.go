package models

import "game/logger"

// テーブル構造
type Profile struct {
	UserID    string  `gorm:"primaryKey;type:varchar(50)" json:"user_id"`        // ユーザID
	RecordID  string  `gorm:"type:varchar(50)" json:"record_id"`                 // 実績ID
	Comment   string  `gorm:"type:varchar(255);default:''" json:"comment"`       // ユーザコメント（デフォルト空白）
	Latitude  float64 `gorm:"type:double;default:0" json:"latitude"`             // 緯度（デフォルト0）
	Longitude float64 `gorm:"type:double;default:0" json:"longitude"`            // 経度（デフォルト0）
	Size      int     `gorm:"default:0" json:"size"`                             // サイズ（デフォルト0）
	RegionID  string  `gorm:"type:varchar(50);default:''" json:"region_id"`      // 地域ID
	SysGame   string  `gorm:"type:varchar(50);default:''" json:"system_game_id"` // システムゲームID
	AdmGame   string  `gorm:"type:varchar(50);default:''" json:"admin_game_id"`  // アドミンゲームID
	Name      string  `gorm:"type:varchar(100);default:''" json:"name"`          //ユーザ名
	Coin      int     `gorm:"default:0" json:"coin"`                             //　所持コイン
}

func (Profile) TableName() string {
	return "Profile"
}

// プロファイルを取得する
func GetProfile(userID string) (*Profile, error) {
	profile := &Profile{}

	// ユーザ情報を取得
	err := dbconn.Where(&Profile{
		UserID: userID,
	}).First(profile).Error

	// エラー処理
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// プロファイルを保存する
func SaveProfile(profile *Profile) error {
	return dbconn.Save(profile).Error
}

// デバッグ用
func DebugProfile() {
	user_id := "userid-e3abf90d-4bcf-4c3b-bbde-37694b1611b3"
	system_game_id := "f5b632cb-707d-f450-eece-f119534b724c"

	//書き込み
	result := dbconn.Save(&Profile{
		UserID:    user_id,
		RecordID:  "第一回優勝",
		Comment:   "よろしくお願いします。",
		Latitude:  35.23,
		Longitude: 135.25,
		Size:      100,
		RegionID:  "関東地方",
		SysGame:   system_game_id,
		AdmGame:   "",
	})

	//エラー処理
	if result.Error != nil {
		logger.PrintErr("プロフィール保存エラー", result.Error)
		return
	}

	logger.PrintErr("プロフィール保存成功")

	//取得コード
	returnData := Sample{}

	//取得する
	result = dbconn.Where(&Sample{
		UserID: user_id,
	}).First(&returnData)

	if result.Error != nil {
		logger.PrintErr("プロフィール取得エラー", result.Error)
		return
	}

	logger.Println("プロフィール取得成功")
}
