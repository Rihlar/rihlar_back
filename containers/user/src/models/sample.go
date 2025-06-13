package models

import "user/logger"

type Sample struct {
	// 独自の主キーを定義
	// primaryKey: 主キーとして設定
	// autoIncrement: 自動増分 (通常は整数型と併用)
	// type:bigint: データ型を指定
	UserID    string    `gorm:"primaryKey;autoIncrement;"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"unique;not null"`
	Age       int       `gorm:"default:0"` // デフォルト値を0に設定
	IsActive  bool      `gorm:"default:true"`
}

func DebugSample() {
	// デバッグ用のコードをここに書く

	userid := "e3abf90d-4bcf-4c3b-bbde-37694b1611b3"

	// 書き込み
	result := dbconn.Save(&Sample{
		UserID:   userid,
		Name:     "aiueo",
		Email:    "test@mattuu.com",
		Age:      20,
		IsActive: false,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("サンプル保存エラー",result.Error)
		return
	}

	logger.Println("サンプル保存成功")

	// 取得コード
	returnData := Sample{}

	// 取得する
	result = dbconn.Where(&Sample{
		UserID:   userid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("サンプル取得エラー",result.Error)
		return
	}

	logger.Println("サンプル取得成功")
}