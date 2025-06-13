package models

import "time"

// テーブル定義
type Circle struct {
	CircleID  string    `gorm:"varchar(36);primaryKey" json:"circlesID"` // サークルID
	GameID    string    `gorm:"varchar(36) not null" json:"gameID"`      // ゲームID
	UserID    string    `gorm:"varchar(36) not null" json:"userID"`      // ユーザーID
	Size      int       `gorm:"not null" json:"size"`                    // サークルサイズ
	Level     int       `gorm:"not null" json:"level"`                   // 防衛レベル
	Latitude  float64   `gorm:"double" json:"latitude"`                  // 緯度
	Longitude float64   `gorm:"double" json:"longitude"`                 // 経度
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`         // 作成時
	ImageID   string    `gorm:"varchar(36)" json:"imageID"`              // イメージID
}

// テーブル名
func (Circle) TableName() string {
	return "Circles"
}
