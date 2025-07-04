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
	// デバッグ用のコードをここに書く

	teamid := "b5fef636-b22e-4057-b1fe-acc7bde6add0"
	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"

	// 書き込み
	result := dbconn.Save(&Team{
		TeamID:    teamid,
		GameID:    gameid,
		Members:   []Member{},
		CreatedAT: time.Now(),
		Points:    515,
	})

	// 使ってないよって言われるので＿使ってます
	_ = dbconn.Save(&Team{
		TeamID:    "4098a6fc-cae8-435d-a24a-48167ec3f3c8",
		GameID:    gameid,
		Members:   []Member{},
		CreatedAT: time.Now(),
		Points:    60,
	})

	_ = dbconn.Save(&Team{
		TeamID:    "e6913e1e-9188-4b21-acfa-aa91ad75d14f",
		GameID:    gameid,
		Members:   []Member{},
		CreatedAT: time.Now(),
		Points:    200,
	})

	_ = dbconn.Save(&Team{
		TeamID:    "608bf57c-427c-423f-8f45-a9f42d337dc9",
		GameID:    "adminGame-7ffcbc90-e8fe-4d9c-8c40-f9f94167dd07",
		Members:   []Member{},
		CreatedAT: time.Now(),
		Points:    0,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チーム保存エラー", result.Error)
		return
	}

	logger.Println("チーム保存成功")

	// 取得コード
	returnData := Team{}

	// 取得する
	result = dbconn.Where(&Team{
		TeamID: teamid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チーム取得エラー", result.Error)
		return
	}

	logger.Println("チーム取得成功")
}
