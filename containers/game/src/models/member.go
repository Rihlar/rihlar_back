package models

import (
	"game/logger"
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

func DebugMember() {
	// デバッグ用のコードをここに書く

	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"
	teamid := "b5fef636-b22e-4057-b1fe-acc7bde6add0"
	userid := "e3abf90d-4bcf-4c3b-bbde-37694b1611b3"

	// 書き込み
	result := dbconn.Save(&Member{
		GameID: gameid,
		TeamID: teamid,
		UserID: userid,
		Points: 500,
	})

	_ = dbconn.Save(&Member{
		GameID: gameid,
		TeamID: teamid,
		UserID: "user-uuid-1",
		Points: 5,
	})

	_ = dbconn.Save(&Member{
		GameID: gameid,
		TeamID: teamid,
		UserID: "user-uuid-2",
		Points: 10,
	})

	_ = dbconn.Save(&Member{
		GameID: gameid,
		TeamID: "4098a6fc-cae8-435d-a24a-48167ec3f3c8",
		UserID: "user-uuid-3",
		Points: 20,
	})

	_ = dbconn.Save(&Member{
		GameID: gameid,
		TeamID: "4098a6fc-cae8-435d-a24a-48167ec3f3c8",
		UserID: "user-uuid-4",
		Points: 40,
	})

	_ = dbconn.Save(&Member{
		GameID: gameid,
		TeamID: "e6913e1e-9188-4b21-acfa-aa91ad75d14f",
		UserID: "user-uuid-5",
		Points: 200,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("メンバー保存エラー", result.Error)
		return
	}

	logger.Println("メンバー保存成功")

	// 取得コード
	returnData := Member{}

	// 取得する
	result = dbconn.Where(&Member{
		UserID: userid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("メンバー取得エラー", result.Error)
		return
	}

	logger.Println("メンバー取得成功")
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
func GetPlaingGames(userUuid string) ([]string, error) {
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
