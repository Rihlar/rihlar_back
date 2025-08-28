package models

import (
	"user/logger"

	"gorm.io/gorm"
)

// テーブル定義
type Member struct {
	GameID string `gorm:"primaryKey;size:50" json:"gameID"` // ゲームID（複合主キー）
	TeamID string `gorm:"not null;size:50" json:"teamID"`   // チームID
	UserID string `gorm:"primaryKey" json:"userID"`         // ユーザーID（複合主キー）
	Points int    `gorm:"not null" json:"points"`           // ポイント
}

// テーブル名
func (Member) TableName() string {
	return "Members"
}

// メンバーを削除する
func (game *Game) DeleteMember(userId string) error {
	return dbconn.Where(&Member{
		UserID: userId,
		GameID: game.GameID,
	}).Delete(&Member{}).Error
}

func DebugMember() {
	logger.Println("メンバー取得成功")
}

// 自身が所属しているチームを取得
func (member *Member) GetTeam() (Team, error) {
	// チームを取得
	var team Team
	err := dbconn.Where(&Team{
		TeamID: member.TeamID,
	}).Find(&team).Error
	return team, err
}


// 現在のランキング取得
func GetMyRanking(userId string, gameId string) (int, error) {

	// ゲームに参加している人間一覧
	var users []Member

	// ゲームIDでフィルタし、得点順に並べて全件取得
	err := dbconn.
		Where("game_id = ?", gameId).
		Order("points DESC").
		Find(&users).Error

	if err != nil {
		return 0, err
	}

	// Go側で DENSE_RANK を計算
	rank := 1        // 現在の順位
	savePoints := -1 // 存在しない得点

	// 人間の数分繰り返す
	for i, s := range users {
		// 一件目は無条件でポイントセーブ
		if i == 0 {
			savePoints = s.Points
		} else if s.Points < savePoints { // 現在の得点が前の得点よりも低ければ、順位を更新
			rank = i + 1
			savePoints = s.Points
		}

		// 現在見ているスコアが自分だったら、順位を返す
		if s.UserID == userId {
			return rank, nil
		}
	}

	return 0, err
}

// 自分の得点取得
func GetMyPoints(userId string, gameId string) (Member, error) {
	var user Member

	// First() は条件に合う最初のレコードを取得し、見つからなければ error が返る
	err := dbconn.Where("user_id = ? AND game_id = ?", userId, gameId).Take(&user).Error
	if err != nil {
		return Member{}, err
	}

	return user, nil
}

// ユーザーの全てのゲームを取得する
func GetJoinGames(userUuid string) ([]string, error) {
	// 結果格納用
	var games []Member

	result := dbconn.Where("user_id = ?", userUuid).Find(&games)
	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return []string{}, nil
	}

	// gameIDだけを抽出
	var gameIds []string
	for _, game := range games {
		gameIds = append(gameIds, game.GameID)
	}
	return gameIds, nil
}

// メンバーを検索する関数
func SearchMember(userId string, gameId string) (Member, error) {
	var member Member

	result := dbconn.Where(Member{
		GameID: gameId,
		UserID: userId,
	}).First(&member)
	if result.Error != nil {
		return Member{}, result.Error
	}

	return member, nil
}

func ExistsMember(userId string, gameId string) (bool, error) {
	member, err := SearchMember(userId, gameId)
	
	// 見つからなければfalse
	if err == gorm.ErrRecordNotFound {
		// 見つからなければfalse
		return false, nil
	}

	// エラー処理
	if err != nil {
		return false, err
	}

	return member.UserID != "", nil
}