package models

import (
	"game/logger"
	"time"
)

// テーブル定義
type Game struct {
	GameID    string    `gorm:"primaryKey;size:50" json:"gameID"`                                                               // ゲームID
	StartTime time.Time `gorm:"DATETIME;not null" json:"startTime"`                                                             // ゲーム開始時間
	EndTime   time.Time `gorm:"DATETIME;not null" json:"endTime"`                                                               // ゲーム終了時間
	Flag      int       `gorm:"not null" json:"flag"`                                                                           // ゲームユニット 0:個人戦、1:チーム戦
	Type      int       `gorm:"not null" json:"type"`                                                                           // ゲームタイプ	0:system、1:admin
	Teams     []Team    `gorm:"foreignKey:GameID;references:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams"` // 参加チーム
	Status    int       `gorm:"not null" json:"status"`                                                                         // ゲームステータス　0:開始前、1:開催中、2:終了済
	RegionID  string    `gorm:"varchar(50)" json:"regionID"`                                                                    // ゲーム開催地域
}

// テーブル名
func (Game) TableName() string {
	return "games"
}

// デバック用
func DebugGame() {
	logger.Println("げーむ取得成功")
}

// メンバーを取得
func (game *Game) GetMemberByUserID(userid string) (Member, error) {
	// 取得する
	returnData := Member{}

	// 取得する
	err := dbconn.Where(&Member{
		UserID: userid,
		GameID: game.GameID,
	}).First(&returnData).Error

	// エラー処理
	if err != nil {
		return Member{}, err
	}

	return returnData, nil
}

// ランキング上位取得
func (game *Game) GetRanking() ([]Team, error) {
	var rankings []Team

	result := dbconn.
		Where("game_id = ?", game.GameID).
		Order("points DESC").
		Find(&rankings)

	if result.Error != nil {
		logger.PrintErr("ランキング上位取得エラー", result.Error)
		return nil, result.Error
	}

	return rankings, nil
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

// ID からゲームを取得
func GetGameByID(gameId string) (Game, error) {
	var game Game

	result := dbconn.Where(&Game{
		GameID: gameId,
	}).First(&game)

	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return Game{}, result.Error
	}

	return game, nil
}
