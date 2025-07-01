package models

import (
	"gcore/logger"
	"time"
)

// テーブル定義
type MovementLog struct {
	UserID    string  `gorm:"primaryKey" json:"userID"`             // ユーザID
	Latitude  float64 `gorm:"double" json:"latitude"`               // 緯度
	Longitude float64 `gorm:"double" json:"longitude"`              // 経度
	Steps     int64   `json:"steps"`                                // 歩数
	GameID    string  `gorm:"primaryKey;varchar(50)" json:"gameID"` // ゲームID
	TimeStamp int64   `gorm:"primaryKey" json:"timeStamp"`          //保存時間
}

// テーブル名
func (MovementLog) TableName() string {
	return "movementLog"
}

// 歩いたログを保存する (緯度経度 歩数)
func (member *Member) SaveMovementLog(Latitude, Longitude float64, Steps int64) error {
	// 歩いた記録をする
	return dbconn.Save(&MovementLog{
		UserID:    member.UserID,
		Latitude:  Latitude,
		Longitude: Longitude,
		Steps:     Steps,
		GameID:    member.GameID,
		TimeStamp: time.Now().Unix(),
	}).Error
}

func (member *Member) GetReportedMovement() ([]MovementLog,error) {
	returnDatas := []MovementLog{}

	// 歩いた記録を取得する
	err := dbconn.Where(&MovementLog{
		UserID: member.UserID,
		GameID: member.GameID,
	}).Order("time_stamp ASC").Find(&returnDatas).Error

	// エラー処理
	if err != nil {
		return []MovementLog{}, err
	}

	return returnDatas, nil
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
		Steps:     10,
		GameID:    gameid,
		TimeStamp: time.Now().Unix(),
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
