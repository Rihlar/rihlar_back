package models

import (
 "game/logger"
	"time"
)

// テーブル定義
type Circle struct {
	CircleID  string    `gorm:"primaryKey" json:"circlesID"`        // サークルID
	GameID    string    `gorm:"varchar(36) not null" json:"gameID"` // ゲームID
	UserID    string    `gorm:"varchar(36) not null" json:"userID"` // ユーザーID
	Size      int       `gorm:"not null" json:"size"`               // サークルサイズ
	Level     int       `gorm:"not null" json:"level"`              // 防衛レベル
	Latitude  float64   `gorm:"double" json:"latitude"`             // 緯度
	Longitude float64   `gorm:"double" json:"longitude"`            // 経度
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`    // 作成時
	ImageID   string    `gorm:"varchar(36)" json:"imageID"`         // イメージID
}

// テーブル名
func (Circle) TableName() string {
	return "Circles"
}

func DebugCircle() {
	// デバッグ用のコードをここに書く

	circleid := "4535e17b-b38c-4449-9902-10861ee3b49b"
	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"
	userid := "e3abf90d-4bcf-4c3b-bbde-37694b1611b3"
	imageid := "76bd1e16-3105-4916-ad6b-7da9554c9601"

	// 書き込み
	result := dbconn.Save(&Circle{
		CircleID:  circleid,
		GameID:    gameid,
		UserID:    userid,
		Size:      1,
		Level:     1,
		Latitude:  34.706414954712386,
		Longitude: 135.50363863029338,
		CreatedAT: time.Time{},
		ImageID:   imageid,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("サークル保存エラー", result.Error)
		return
	}

	logger.Println("サークル保存成功")

	// 取得コード
	returnData := Circle{}

	// 取得する
	result = dbconn.Where(&Circle{
		CircleID: circleid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("サークル取得エラー", result.Error)
		return
	}

	logger.Println("サークル取得成功")
}


// 円の詳細取得
func GetCircleDeteile(circleId string) (Circle, error) {
	var circleDeteile Circle

	result := dbconn.Where("circle_id = ?", circleId).Take(&circleDeteile)
	if result.Error != nil {
		logger.PrintErr("円詳細取得エラー", result.Error)
		return Circle{}, nil
	}

	return circleDeteile, nil
}