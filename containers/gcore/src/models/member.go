package models

import "gcore/logger"

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

// レベル2の円を取得する関数
func (member *Member) GetLevelTwoCircles() ([]Circle, error) {
	// 取得する
	returnDatas := []Circle{}

	// 取得する
	err := dbconn.Where(&Circle{
		UserID: member.UserID,
		Level:   2,
	}).Find(&returnDatas).Error

	// エラー処理
	if err != nil {
		return []Circle{}, err
	}

	return returnDatas, nil
}

// 点数を更新する関数
func (member *Member) GetOwnerdChunks() ([]GameChunk, error) {
	// チャンクを取得する
	returnChunks := []GameChunk{}

	// チャンクを検索する
	err := dbconn.Where(&GameChunk{
		OwnerID: member.UserID,
		GameID:  member.GameID,
	}).Find(&returnChunks).Error

	// エラー処理
	if err != nil {
		return []GameChunk{}, err
	}

	return returnChunks, nil
}

// ポイントを更新する関数
func (member *Member) UpdatePoints(points int) error {
	// 更新する
	member.Points = points

	// 更新する
	return dbconn.Model(member).Update("points", points).Error
}

// 所有中のチャンクをポイントに反映する関数
func (member *Member) ReflectPoints() error {
	// 所有中のチャンクを取得する
	chunks, err := member.GetOwnerdChunks()

	// エラー処理
	if err != nil {
		return err
	}

	// 更新する
	err = member.UpdatePoints(len(chunks))

	// エラー処理
	if err != nil {
		return err
	}

	return nil
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
