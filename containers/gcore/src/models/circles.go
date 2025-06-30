package models

import (
	"gcore/logger"
	"time"
)

// テーブル定義
type Circle struct {
	CircleID  string    `gorm:"primaryKey" json:"circlesID"`        // サークルID
	GameID    string    `gorm:"varchar(50) not null" json:"gameID"` // ゲームID
	UserID    string    `gorm:"varchar(50) not null" json:"userID"` // ユーザーID
	Size      int       `gorm:"not null" json:"size"`               // サークルサイズ
	Level     int       `gorm:"not null" json:"level"`              // 防衛レベル
	Latitude  float64   `gorm:"double" json:"latitude"`             // 緯度
	Longitude float64   `gorm:"double" json:"longitude"`            // 経度
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`    // 作成時
	ImageID   string    `gorm:"varchar(50)" json:"imageID"`         // イメージID
	Steps     int64     `json:"steps"`                              // 歩数
}

// テーブル名
func (Circle) TableName() string {
	return "Circles"
}

// 円を取得する関数
func GetCircle(circleid string) (Circle, error) {
	var circle Circle
	return circle, dbconn.Where(&Circle{
		CircleID: circleid,
	}).First(&circle).Error
}

// メンバーに属する円の一覧取得
func (member *Member) GetCircles() ([]Circle, error) {
	// 取得
	returnData := []Circle{}

	// 取得
	return returnData, dbconn.Where(&Circle{
		UserID: member.UserID,
	}).Find(&returnData).Error
}

// 円を作成する
func (member *Member) CreateCircle(circle *Circle) error {
	// 値を設定
	circle.GameID = member.GameID
	circle.UserID = member.UserID

	// 作成時の時間を設定
	circle.CreatedAT = time.Now()

	// 円を作成
	err := dbconn.Save(circle).Error

	// エラー処理
	if err != nil {
		return err
	}

	// 円のレベルを更新 (関連するチャンクを更新するため)
	return circle.ChangeLevel(circle.Level)
}

// 円のレベルを変更する
func (circle *Circle) ChangeLevel(level int) error {
	// 変更
	circle.Level = level

	// ゲームを取得する
	game, err := GetGame(circle.GameID)

	// エラー処理
	if err != nil {
		return err
	}

	// 円形にチャンクを取得する
	chunks, err := game.GetCircleChunkByLatLon(circle.Latitude, circle.Longitude,float64(circle.Size))

	// エラー処理
	if err != nil {
		return err
	}

	for _, chunk := range chunks {
		// チャンクのレベルがこの円のレベル以下なら
		if chunk.Level <= circle.Level {
			// チャンクを更新する
			if err := chunk.ChangeLevel(circle.Level); err != nil {
				logger.PrintErr(err)
				return err
			}

			// オーナーも変更する
			if err := chunk.ChangeOwner(circle.UserID); err != nil {
				logger.PrintErr(err)
				return err
			}
		}
	}

	// 更新
	return dbconn.Model(circle).Update("level", level).Error
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
