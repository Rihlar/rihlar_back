package models

import "time"

// テーブル定義
type Team struct {
	TeamID    string    `gorm:"varchar(36);primaryKey" json:"teamID"`                                           // チームID
	GameID    string    `gorm:"varchar(36) not null" json:"gameID"`                                             // ゲームID
	Members   []Member  `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"` //　チームメンバー
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`                                                // ゲーム作成時
	Points    int       `gorm:"not null" json:"points"`                                                         // チーム合計ポイント
}

// テーブル名
func (Team) TableName() string {
	return "Teams"
}
