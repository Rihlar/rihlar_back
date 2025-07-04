package models

import "gcore/logger"

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
}

func GetProfile(userID string) (Profile, error) {
	var profile Profile

	// ユーザ情報を取得
	err := dbconn.Where(&Profile{
		UserID: userID,
	}).First(&profile).Error

	// エラー処理
	if err != nil {
		logger.PrintErr("ユーザ情報取得エラー", err)
		return Profile{}, err
	}

	return profile, nil
}

func SaveProfile(profile Profile) error {
	return dbconn.Save(&profile).Error
}

func (Profile) TableName() string {
	return "Profile"
}

// デバッグ用
func DebugProfile() {
	// プロファイルを作成する
	CreateTestProfiles()

	logger.Println("プロフィール取得成功")
}

func CreateTestProfiles() {
	for index, userid := range UserIDs {
		CreateProfile(Profile{
			UserID:    userid,
			RecordID:  "",
			Comment:   "",
			Latitude:  0,
			Longitude: 0,
			Size:      0,
			RegionID:  "",
			SysGame:   SysGameIDs[index],
			AdmGame:   AdminGameId1,
			Name:      AllUserNames[index],
		})
	}
}

func CreateProfile(data Profile) error {
	return dbconn.Create(&data).Error
}