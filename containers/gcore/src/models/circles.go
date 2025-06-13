package models

import "time"

// テーブル定義
type Circle struct {
	CircleID string    `gorm:"varchar(36);primaryKey" json:"circlesID"`    // ゲームID
	GameID    string    `gorm:"varchar(36) not null" json:"gameID"` // ゲーム開始時間
	UserID    string    `gorm:"varchar(36) not null" json:"userID"` // ゲーム終了時間
	Size      int       `gorm:"int not null" json:"size"`           // ゲームユニット 0:個人戦、1:チーム戦
	Level     int       `gorm:"int not null" json:"level"`          // ゲームタイプ	0:system、1:admin
	Latitude  float64   `gorm:"double" json:"latitude"`             // 参加チーム
	Longitude float64   `gorm:"double" json:"longitude"`            // ゲームステータス　0:開始前、1:開催中、2:終了済
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`            // ゲーム開催地域
	ImageID   string    `gorm:"varchar(36)" json:"imageID"`         // ゲーム開催地域
}

// テーブル名
func (Circle) TableName() string {
	return "Circles"
}
