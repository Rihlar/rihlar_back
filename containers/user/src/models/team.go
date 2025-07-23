package models

import (
	"time"
	"user/logger"
)

// テーブル定義
type Team struct {
	TeamID    string    `gorm:"primaryKey;size:50" json:"teamID"`                                                                 // チームID
	GameID    string    `gorm:"not null;size:50" json:"gameID"`                                                                   // ゲームID
	Members   []Member  `gorm:"foreignKey:TeamID,GameID;references:TeamID,GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"members"` //　チームメンバー
	CreatedAT time.Time `gorm:"autoCreateTime" json:"createdAT"`                                                                  // ゲーム作成時
	Points    int       `gorm:"not null" json:"points"`                                                                           // チーム合計ポイント
}

// チームを作成する
func CreateTeam(team Team) error {
	return dbconn.Create(&team).Error
}

// TODO デバッグ用 メンバーを追加する
func (team *Team) AddMember(member Member) error {
	return dbconn.Model(team).Association("Members").Append(&member)
}

// チームを取得する
func (game *Game) GetTeam(teamID string) (Team, error) {
	// 取得する
	returnData := Team{}

	// 取得
	err := dbconn.Where(&Team{
		GameID: game.GameID,
		TeamID: teamID,
	}).First(&returnData).Error

	// エラー処理
	if err != nil {
		return Team{}, err
	}

	return returnData, nil
}

// テーブル名
func (Team) TableName() string {
	return "Teams"
}

// チームを削除する
func (game *Game) DeleteTeam(teamID string) error {
	return dbconn.Where(&Team{
		TeamID: teamID,
		GameID: game.GameID,
	}).Delete(&Team{}).Error
}

// 自分のチームのランキングを取得
func (team *Team) GetRank() (int, error) {
	// 取得する
	teams := []Team{}

	// チームをソートして取得
	result := dbconn.Where(&Team{
		GameID: team.GameID,
	}).Order("points desc").Find(&teams)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("ランキング上位取得エラー", result.Error)
		return 0, result.Error
	}

	// 自身のチームが来るまでチームを回す
	for index, checkTeam := range teams {
		if checkTeam.TeamID == team.TeamID {
			// 自身のチームが来た時にindexを返す
			return index + 1, nil
		}
	}

	return 0, nil
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