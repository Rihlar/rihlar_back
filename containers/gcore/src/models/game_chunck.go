package models

// テーブル定義
type GameChunck struct {
	ChunkID string `gorm:"varchar(36);primaryKey" json:"chunkID"` // チャンクID
	GameID  string `gorm:"varchar(36)" json:"gameID"`             // ゲームID
	ImageID string `gorm:"varchar(36)" json:"imageID"`            // イメージID
	OwnerID string `gorm:"varchar(36)" json:"ownerID"`            // オーナーID
	Level   int    `gorm:"not null" json:"level"`                 // 防衛レベル
}

// テーブル名
func (GameChunck) TableName() string {
	return "game_chunks"
}

