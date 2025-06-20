package models

import (
	"gcore/logger"
	"time"
)

// テーブル定義
type Game struct {
	GameID    string    `gorm:"primaryKey;size:36" json:"gameID"`                                                     // ゲームID
	StartTime time.Time `gorm:"DATETIME;not null" json:"startTime"`                                           // ゲーム開始時間
	EndTime   time.Time `gorm:"DATETIME;not null" json:"endTime"`                                             // ゲーム終了時間
	Flag      int       `gorm:"not null" json:"flag"`                                                         // ゲームユニット 0:個人戦、1:チーム戦
	Type      int       `gorm:"not null" json:"type"`                                                         // ゲームタイプ	0:system、1:admin
	Teams     []Team    `gorm:"foreignKey:GameID;references:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams"` // 参加チーム
	Status    int       `gorm:"not null" json:"status"`                                                       // ゲームステータス　0:開始前、1:開催中、2:終了済
	RegionID  string    `gorm:"varchar(36)" json:"regionID"`                                                  // ゲーム開催地域
}

// テーブル名
func (Game) TableName() string {
	return "games"
}

// ゲームを取得するエンドポイント
func GetGame(gameid string) (Game, error) {
	var game Game

	// 取得する
	result := dbconn.Where(&Game{
		GameID: gameid,
	}).Find(&game)

	return game, result.Error
}

// ゲームを保存するエンドポイント
func SaveGame(game Game) error {
	return dbconn.Save(&game).Error
}

// デバック用
func DebugGame() {
	// デバッグ用のコードをここに書く

	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"
	regionid := "f6b4e846-1e99-45a1-a7a7-1858a9f94d28" // kansai

	// 書き込み
	result := dbconn.Save(&Game{
		GameID:    gameid,
		StartTime: time.Now().AddDate(0, 0, 1),
		EndTime:   time.Now().AddDate(0, 0, 5),
		Flag:      0,
		Type:      1,
		Teams:     []Team{},
		Status:    0,
		RegionID:  regionid,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("ゲーム保存エラー", result.Error)
		return
	}

	logger.Println("ゲーム保存成功")

	// 取得コード
	returnData := Game{}

	// 取得する
	result = dbconn.Where(&Game{
		GameID: gameid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return
	}

	logger.Println("げーむ取得成功")
}

