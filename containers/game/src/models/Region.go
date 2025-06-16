package models

import "game/logger"

// テーブル定義
type Region struct {
	RegionID   string  `gorm:"primaryKey" json:"RegionID"`   // リージョンID
	RegionName string `gorm:"not null" json:"regionName"`	// リージョン名
}

// テーブル名
func (Region) TableName() string {
	return "regions"
}

func DebugRegion() {
	// デバッグ用のコードをここに書く

	regionid := "f6b4e846-1e99-45a1-a7a7-1858a9f94d28" // kansai

	// 書き込み
	result := dbconn.Save(&Region{
		RegionID:   regionid,
		RegionName: "kansai",
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("リージョン保存エラー", result.Error)
		return
	}

	logger.Println("リージョン保存成功")

	// 取得コード
	returnData := Region{}

	// 取得する
	result = dbconn.Where(&Region{
		RegionID: regionid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("リージョン取得エラー", result.Error)
		return
	}

	logger.Println("リージョン取得成功")
}
