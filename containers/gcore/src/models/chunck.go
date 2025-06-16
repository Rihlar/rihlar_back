package models

import "gcore/logger"

// テーブル定義
type Chunk struct {
	ChunkID   string  `gorm:"primaryKey" json:"chunkID"`   // チャンクID
	Latitude  float64 `gorm:"double" json:"latitude"`      // 緯度
	Longitude float64 `gorm:"double" json:"longitude"`     // 経度
	RegionID  string  `gorm:"varchar(36)" json:"regionID"` // ゲーム開催地域
}

// テーブル名
func (Chunk) TableName() string {
	return "chunks"
}

func DebugChunk() {
	// デバッグ用のコードをここに書く

	chunkid := "3325d4ee-ef32-42a3-91d1-33d3582dffc2"
	regionid := "f6b4e846-1e99-45a1-a7a7-1858a9f94d28" // kansai

	// 書き込み
	result := dbconn.Save(&Chunk{
		ChunkID:   chunkid,
		Latitude:  0,
		Longitude: 0,
		RegionID:  regionid,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チャンク保存エラー", result.Error)
		return
	}

	logger.Println("チャンク保存成功")

	// 取得コード
	returnData := Chunk{}

	// 取得する
	result = dbconn.Where(&Chunk{
		ChunkID: chunkid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チャンク取得エラー", result.Error)
		return
	}

	logger.Println("チャンク取得成功")
}
