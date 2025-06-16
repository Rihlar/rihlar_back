package models

import "game/logger"

// テーブル定義
type BaseChunk struct {
	ChunkID   string  `gorm:"primaryKey" json:"chunkID"`   // チャンクID
	Latitude  float64 `gorm:"double" json:"latitude"`      // 緯度
	Longitude float64 `gorm:"double" json:"longitude"`     // 経度
	RegionID  string  `gorm:"varchar(36)" json:"regionID"` // ゲーム開催地域
}

// テーブル名
func (BaseChunk) TableName() string {
	return "chunks"
}

func DebugBaseChunk() {
	// デバッグ用のコードをここに書く

	chunkid := "3325d4ee-ef32-42a3-91d1-33d3582dffc2"
	regionid := "f6b4e846-1e99-45a1-a7a7-1858a9f94d28" // kansai

	// 書き込み
	result := dbconn.Save(&BaseChunk{
		ChunkID:   chunkid,
		Latitude:  0,
		Longitude: 0,
		RegionID:  regionid,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("ベースチャンク保存エラー", result.Error)
		return
	}

	logger.Println("ベースチャンク保存成功")

	// 取得コード
	returnData := BaseChunk{}

	// 取得する
	result = dbconn.Where(&BaseChunk{
		ChunkID: chunkid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("ベースチャンク取得エラー", result.Error)
		return
	}

	logger.Println("ベースチャンク取得成功")
}
