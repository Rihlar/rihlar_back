package models

// テーブル定義
type GameChunk struct {
	ChunkID  string  `gorm:"primaryKey" json:"chunkID"`  // チャンクID
	GameID   string  `gorm:"varchar(50)" json:"gameID"`  // ゲームID
	ImageID  string  `gorm:"varchar(50)" json:"imageID"` // イメージID
	OwnerID  string  `gorm:"varchar(50)" json:"ownerID"` // オーナーID
	StartLat float64 `gorm:"not null" json:"startLat"`   // 開始緯度
	StartLon float64 `gorm:"not null" json:"startLon"`   // 開始経度
	EndLat   float64 `gorm:"not null" json:"endLat"`     // 終了緯度
	EndLon   float64 `gorm:"not null" json:"endLon"`     // 終了経度
	Level    int     `gorm:"not null" json:"level"`      // 防衛レベル
	
}

// テーブル名
func (GameChunk) TableName() string {
	return "game_chunks"
}

func DebugGameChunk() {
}