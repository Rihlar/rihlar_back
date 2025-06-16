package models

import "gcore/logger"

// テーブル定義
type Member struct {
	GameID string `gorm:"primaryKey;size:36" json:"gameID"` // ゲームID（複合主キー）
	TeamID string `gorm:"not null;size:36" json:"teamID"`   // チームID
	UserID string `gorm:"primaryKey" json:"userID"` // ユーザーID（複合主キー）
	Points int    `gorm:"not null" json:"points"`   // ポイント
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
		Points: 0,
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
