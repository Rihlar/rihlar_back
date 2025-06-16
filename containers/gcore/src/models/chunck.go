package models

import "gcore/logger"

// テーブル定義
type Chunck struct {
	ChunkID   string  `gorm:"varchar(36);primaryKey" json:"chunkID"` // チャンクID
	Latitude  float64 `gorm:"double" json:"latitude"`                // 緯度
	Longitude float64 `gorm:"double" json:"longitude"`               // 経度
	RegionID  string  `gorm:"varchar(36)" json:"regionID"`           // ゲーム開催地域
}

// テーブル名
func (Chunck) TableName() string {
	return "chunks"
}

func DebugChunck() {
	// デバッグ用のコードをここに書く

	chunkid := "3325d4ee-ef32-42a3-91d1-33d3582dffc2"

	// 書き込み
	result := dbconn.Save(&Chunck{
		ChunkID:   chunkid,
		Latitude:  0,
		Longitude: 0,
		RegionID:  "",
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("サークル保存エラー",result.Error)
		return
	}

	logger.Println("サークル保存成功")

	// 取得コード
	returnData := Chunck{}

	// 取得する
	result = dbconn.Where(&Chunck{
		ChunkID:   chunkid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("サークル取得エラー",result.Error)
		return
	}

	logger.Println("サークル取得成功")
}