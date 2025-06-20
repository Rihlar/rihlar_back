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

// 点数を更新する
func (member *Member) UpdatePoints(point int) error {
	// 点数を更新
	member.Points = point

	// アップデートを実行する
	result := dbconn.Model(member).Update("points", member.Points)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("メンバー更新エラー", result.Error)
		return result.Error
	}

	return nil
}

// メンバーを取得する
func GetMemberByUserID(userid string) (Member, error) {
	var member Member

	// 取得する
	result := dbconn.Where(&Member{
		UserID: userid,
	}).First(&member)

	return member, result.Error
}

// メンバーを保存する (上書き)
func SaveMember(member Member) error {
	return dbconn.Save(&member).Error
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
