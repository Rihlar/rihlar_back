package models

import "game/logger"

// テーブル定義
type MovementLog struct {
	UserID    string  `gorm:"primaryKey" json:"userID"` // ユーザID
	Latitude  float64 `gorm:"double" json:"latitude"`     // 緯度
	Longitude float64 `gorm:"double" json:"longitude"`    // 経度
	GameID    string  `gorm:"varchar(36)" json:"gameID"`  // ゲームID
}

// テーブル名
func (MovementLog) TableName() string {
	return "movementLog"
}

func DebugMovementLog() {
	// デバッグ用のコードをここに書く

	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"
	userid := "e3abf90d-4bcf-4c3b-bbde-37694b1611b3"

	// 書き込み
	result := dbconn.Save(&MovementLog{
		UserID:    userid,
		Latitude:  0,
		Longitude: 0,
		GameID:    gameid,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("movement log保存エラー", result.Error)
		return
	}

	logger.Println("movement log保存成功")

	// 取得コード
	returnData := MovementLog{}

	// 取得する
	result = dbconn.Where(&MovementLog{
		UserID: userid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("movementログ取得エラー", result.Error)
		return
	}

	logger.Println("movementログ取得成功")
}


