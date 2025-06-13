package models

import "time"

// テーブル定義
type Game struct {
	GameID    string    `gorm:"varchar(36);primaryKey" json:"gameID"`                                         // ゲームID
	StartTime time.Time `gorm:"DATETIME;not null" json:"startTime"`                                           // ゲーム開始時間
	EndTime   time.Time `gorm:"DATETIME;not null" json:"endTime"`                                             // ゲーム終了時間
	Flag      int       `gorm:"not null" json:"flag"`                                                         // ゲームユニット 0:個人戦、1:チーム戦
	Type      int       `gorm:"not null" json:"type"`                                                         // ゲームタイプ	0:system、1:admin
	Teams     []Team    `gorm:"foreignKey:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams"` // 参加チーム
	Status    int       `gorm:"not null" json:"status"`                                                       // ゲームステータス　0:開始前、1:開催中、2:終了済
	RegionID  string    `gorm:"varchar(36)" json:"regionID"`                                                  // ゲーム開催地域
}

// テーブル名
func (Game) TableName() string {
	return "games"
}
