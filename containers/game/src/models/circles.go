package models

import (
	"game/logger"
	"time"
)

// テーブル定義
type Circle struct {
	CircleID  string    `gorm:"primaryKey" json:"circlesID"`        // サークルID
	GameID    string    `gorm:"varchar(50) not null" json:"gameID"` // ゲームID
	TeamID    string    `gorm:"varchar(50) not null" json:"teamID"` // チームID
	UserID    string    `gorm:"varchar(50) not null" json:"userID"` // ユーザーID
	Size      int       `gorm:"not null" json:"size"`               // サークルサイズ
	Level     int       `gorm:"not null" json:"level"`              // 防衛レベル
	Latitude  float64   `gorm:"double" json:"latitude"`             // 緯度
	Longitude float64   `gorm:"double" json:"longitude"`            // 経度
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`    // 作成時
	ImageID   string    `gorm:"varchar(50)" json:"imageID"`         // イメージID
	Steps     int64     `json:"steps"`                              // 歩数
	Theme     string    `json:"theme"`								// 円のテーマ
}

// テーブル名
func (Circle) TableName() string {
	return "Circles"
}

func DebugCircle() {
	logger.Println("サークル取得成功")
}


// 円の詳細取得
func GetCircleDeteile(circleId string) (Circle, error) {
	var circleDeteile Circle

	result := dbconn.Where("circle_id = ?", circleId).Take(&circleDeteile)
	if result.Error != nil {
		logger.PrintErr("円詳細取得エラー", result.Error)
		return Circle{}, nil
	}

	return circleDeteile, nil
}

// チームのサークルを取得
func GetCircleByTeamId(teamId string) ([]Circle, error) {
	var circles []Circle

	// 取得
	err := dbconn.Where(&Circle{
		TeamID: teamId,
	}).Find(&circles).Error

	// エラー処理
	if err != nil {
		logger.PrintErr("サークル取得エラー", err)
		return []Circle{}, err
	}

	return circles, nil
}

// ゲームに属する円を取得する
func (game Game) GetCircles() ([]Circle, error) {
	var circles []Circle

	// 取得
	err := dbconn.Where(&Circle{
		GameID: game.GameID,
	}).Find(&circles).Error

	// エラー処理
	if err != nil {
		logger.PrintErr("サークル取得エラー", err)
		return []Circle{}, err
	}

	return circles, nil
}