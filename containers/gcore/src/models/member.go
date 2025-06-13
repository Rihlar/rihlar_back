package models

// テーブル定義
type Member struct {
	GameID string `gorm:"varchar(36);primaryKey" json:"gameID"` // ゲームID（複合主キー）
	TeamID string `gorm:"varchar(36);not null" json:"teamID"`   // チームID
	UserID string `gorm:"varchar(36);primaryKey" json:"userID"` // ユーザーID（複合主キー）
	Points int    `gorm:"not null" json:"points"`               // ポイント
}

// テーブル名
func (Member) TableName() string {
	return "Members"
}
