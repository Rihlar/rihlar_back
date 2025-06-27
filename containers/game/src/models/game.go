package models

import (
	"game/logger"
	"time"
)

// テーブル定義
type Game struct {
	GameID    string    `gorm:"primaryKey;size:36" json:"gameID"`                                                               // ゲームID
	StartTime time.Time `gorm:"DATETIME;not null" json:"startTime"`                                                             // ゲーム開始時間
	EndTime   time.Time `gorm:"DATETIME;not null" json:"endTime"`                                                               // ゲーム終了時間
	Flag      int       `gorm:"not null" json:"flag"`                                                                           // ゲームユニット 0:個人戦、1:チーム戦
	Type      int       `gorm:"not null" json:"type"`                                                                           // ゲームタイプ	0:system、1:admin
	Teams     []Team    `gorm:"foreignKey:GameID;references:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams"` // 参加チーム
	Status    int       `gorm:"not null" json:"status"`                                                                         // ゲームステータス　0:開始前、1:開催中、2:終了済
	RegionID  string    `gorm:"varchar(36)" json:"regionID"`                                                                    // ゲーム開催地域
}

// テーブル名
func (Game) TableName() string {
	return "games"
}

// デバック用
func DebugGame() {
	// デバッグ用のコードをここに書く

	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"
	regionid := "f6b4e846-1e99-45a1-a7a7-1858a9f94d28" // kansai

	// 書き込み　開始中ゲーム
	result := dbconn.Save(&Game{
		GameID:    gameid,
		StartTime: time.Now().AddDate(0, 0, 1),
		EndTime:   time.Now().AddDate(0, 0, 5),
		Flag:      0,
		Type:      1,
		Teams:     []Team{},
		Status:    1,
		RegionID:  regionid,
	})

	_ = dbconn.Save(&Game{
		GameID:    "a7510bcb-d5b8-414b-84ef-d4c663452e43",
		StartTime: time.Now().AddDate(0, 0, 1),
		EndTime:   time.Now().AddDate(0, 0, 3),
		Flag:      0,
		Type:      1,
		Teams:     []Team{},
		Status:    2,
		RegionID:  regionid,
	})

	_ = dbconn.Save(&Game{
		GameID:    "3b7f30f3-5abe-41e6-8ea6-fc8fdc45a51c",
		StartTime: time.Now().AddDate(0, 0, 4),
		EndTime:   time.Now().AddDate(0, 0, 6),
		Flag:      0,
		Type:      1,
		Teams:     []Team{},
		Status:    2,
		RegionID:  regionid,
	})

	// systemゲーム
	_ = dbconn.Save(&Game{
		GameID:    "ec2d55e6-d89b-4821-8933-ef2d4bb18703",
		StartTime: time.Now().AddDate(0, 0, 7),
		EndTime:   time.Now().AddDate(0, 0, 9),
		Flag:      0,
		Type:      0,
		Teams:     []Team{},
		Status:    1,
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

// ゲームの詳細取得
func GetGame(gameId []string) ([]Game, error) {
	// 結果格納用
	var games []Game

	result := dbconn.Where("game_id IN ?", gameId).Find(&games)
	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return []Game{}, nil
	}

	return games, nil
}

// 開催中のゲーム取得
func GetGameHolding(gameId []string) ([]Game, error) {
	// 結果格納用
	var games []Game

	result := dbconn.Where("game_id IN ?", gameId).Where("status = ?", 1).Find(&games)
	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return []Game{}, nil
	}

	return games, nil
}

// 終了済みゲーム一覧取得
func GetEndGames(gameId []string) ([]Game, error) {
	// 結果格納用
	var games []Game

	// statusが2で絞る
	result := dbconn.Where("game_id IN ?", gameId).Where("status = ?", 2).Find(&games)
	if result.Error != nil {
		logger.PrintErr("ゲーム取得エラー", result.Error)
		return []Game{}, nil
	}

	return games, nil
}
