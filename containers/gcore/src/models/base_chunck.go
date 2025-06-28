package models

// テーブル定義
type BaseChunk struct {
	BaseChunkID string  `gorm:"primaryKey" json:"chunkID"`   // チャンクID
	Latitude    float64 `gorm:"double" json:"latitude"`      // 緯度
	Longitude   float64 `gorm:"double" json:"longitude"`     // 経度
	RegionID    string  `gorm:"varchar(50)" json:"regionID"` // ゲーム開催地域
}

// テーブル名
func (BaseChunk) TableName() string {
	return "chunks"
}

func DebugBaseChunk() {
}
