package models

import (
	"game/logger"
	"time"
)

// テーブル定義
type Team struct {
	TeamID    string    `gorm:"primaryKey;size:50" json:"teamID"`                                                                 // チームID
	GameID    string    `gorm:"not null;size:50" json:"gameID"`                                                                   // ゲームID
	Members   []Member  `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"` //　チームメンバー
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`                                                                  // ゲーム作成時
	Points    int       `gorm:"not null" json:"points"`                                                                           // チーム合計ポイント
}

// テーブル名
func (Team) TableName() string {
	return "Teams"
}

func (game *Game) GetTeamByUserID(UserID string) (Team, error) {
	// 取得する
	returnData := Team{}

	// メンバーを取得
	member, err := game.GetMemberByUserID(UserID)

	// エラー処理
	if err != nil {
		return Team{}, err
	}

	// 取得する
	err = dbconn.Where(&Team{
		TeamID: member.TeamID,
	}).First(&returnData).Error

	// エラー処理
	if err != nil {
		return Team{}, err
	}

	return returnData, nil
}

func DebugTeam() {
	logger.Println("チーム取得成功")
}
