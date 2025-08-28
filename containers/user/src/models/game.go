package models

import (
	"time"
	"user/logger"
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

// TODO デバッグ用 チームを追加する
func (game *Game) AddTeam(team Team) error {
	return dbconn.Model(game).Association("Teams").Append(&team)
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

// ランキング上位取得
func (game *Game) GetRanking(maxRank int) ([]Team, error) {
	var rankings []Team

	result := dbconn.Debug().
		Where(Team{
			GameID: game.GameID,
		}).
		Order("points DESC").
		Limit(maxRank).
		Find(&rankings)

	if result.Error != nil {
		logger.PrintErr("ランキング上位取得エラー", result.Error)
		return nil, result.Error
	}

	return rankings, nil
}

func CreateGame(game Game) error {
	return dbconn.Create(&game).Error
}

// 一人のゲーム追加をデバッグする
func DebugAddMember(gameID string, teamID string, userID string) error {
	// ゲームを取得する
	game, err := GetGame(gameID)

	// エラー処理
	if err != nil {
		logger.PrintErr("ゲーム取得エラー", err)
		return err
	}

	// ゲームにチームを追加する
	err = game.AddTeam(Team{
		TeamID: teamID,
		Points: 0,
	})

	// エラー処理
	if err != nil {
		return err
	}

	// チームを取得
	team, err := game.GetTeam(teamID)

	// エラー処理
	if err != nil {
		return err
	}

	// チームにメンバーを追加する
	err = team.AddMember(Member{
		UserID: userID,
		Points: 0,
	})

	// エラー処理
	if err != nil {
		logger.PrintErr("メンバー追加エラー", err)
		return err
	}

	return nil
}
