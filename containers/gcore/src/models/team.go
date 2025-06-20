package models

import (
	"gcore/logger"
	"time"
)

// テーブル定義
type Team struct {
	TeamID    string    `gorm:"primaryKey;size:36" json:"teamID"`                                           // チームID
	GameID    string    `gorm:"not null;size:36" json:"gameID"`                                             // ゲームID
	Members   []Member  `gorm:"foreignKey:TeamID;references:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"` //　チームメンバー
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`                                                // ゲーム作成時
	Points    int       `gorm:"not null" json:"points"`                                                         // チーム合計ポイント
}

// テーブル名
func (Team) TableName() string {
	return "Teams"
}

// チームの点数を更新する
func (team *Team) UpdatePoints(point int) error {
	// チームの点数を更新
	team.Points = point

	// アップデートを実行する
	result := dbconn.Model(team).Update("points", team.Points)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チーム更新エラー", result.Error)
		return result.Error
	}

	return nil
}

// チームを取得する処理
func GetTeam(teamID string) (Team, error) {
	returnData := Team{}

	// 取得する
	result := dbconn.Where(&Team{
		TeamID: teamID,
	}).First(&returnData)

	return returnData, result.Error
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
		CreatedAT: time.Time{},
		Points:    0,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チーム保存エラー",result.Error)
		return
	}

	logger.Println("チーム保存成功")

	// 取得コード
	returnData := Team{}

	// 取得する
	result = dbconn.Where(&Team{
		TeamID:   teamid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("チーム取得エラー",result.Error)
		return
	}

	logger.Println("チーム取得成功")
}