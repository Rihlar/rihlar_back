package models

import (
	"gcore/logger"
	"time"
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

// デバック用
func DebugGame() {
	// デバッグ用のコードをここに書く

	// admin のゲームを作成する
	debugAdminGame()

	// system のゲームを作成する
	debugSystemGame()

	// ゲームを作成したのでデバッグ用のユーザーをゲームに追加する
	debugGameUser()

	// メンバーを追加するをテストする
	DebugAddMember(admin_game_id, teamID, user_id)
	DebugAddMember(admin_game_id, teamID2, user_id2)
	DebugAddMember(system_game_id, sysTeamID, user_id)
	DebugAddMember(system_game_id2, sysTeamID, user_id2)
}

// 管理者が作成したゲームのデバッグをする
func debugAdminGame() {
	// ゲームを作成する
	err := SaveGame(Game{
		GameID:    admin_game_id,
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0, 0, 5),
		Flag:      0,
		Type:      1,
		Teams:     []Team{},
		Status:    0,
		RegionID:  regionid,
	})

	if err != nil {
		logger.PrintErr("ゲーム作成エラー", err)
		return
	}

	// ゲームを取得する
	game, err := GetGame(admin_game_id)

	// エラー処理
	if err != nil {
		logger.PrintErr("ゲーム取得エラー", err)
		return
	}

	// 関西リージョン取得
	region, err := GetRegionByID(regionid)

	// エラー処理
	if err != nil {
		logger.PrintErr("リージョン取得エラー", err)
		return
	}

	// チャンクを作成する
	err = game.FillRegion(region)

	// エラー処理
	if err != nil {
		logger.PrintErr("チャンク作成エラー", err)
		return
	}

	logger.Println("ゲーム作成成功")
}

// システムが作成したゲームのデバッグをする
func debugSystemGame() {
	// ゲームを作成する
	err := SaveGame(Game{
		GameID:    system_game_id,
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0, 0, 5),
		Flag:      0,
		Type:      0,
		Teams:     []Team{},
		Status:    1,
		RegionID:  regionid,
	})

	if err != nil {
		logger.PrintErr("ゲーム作成エラー", err)
		return
	}

	logger.Println("ゲーム作成成功")

	// 二つ目のシステムゲームを作成
	err = SaveGame(Game{
		GameID:    system_game_id2,
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0, 0, 5),
		Flag:      0,
		Type:      0,
		Teams:     []Team{},
		Status:    1,
		RegionID:  regionid,
	})

	if err != nil {
		logger.PrintErr("ゲーム作成エラー", err)
		return
	}
}

// ゲームにユーザーを追加する
func debugGameUser() {
	// ユーザーID

	// プロファイルを取得する
	profile, err := GetProfile(user_id)

	// エラー処理
	if err != nil {
		logger.PrintErr("プロファイル取得エラー", err)
		return
	}

	// プロファイルを更新する
	profile.SysGame = system_game_id
	profile.AdmGame = admin_game_id

	// プロファイルを保存する
	err = SaveProfile(profile)

	// エラー処理
	if err != nil {
		logger.PrintErr("プロファイル保存エラー", err)
		return
	}

	// 二人目のプロファイルを更新する
	profile2, err := GetProfile(user_id2)

	// エラー処理
	if err != nil {
		logger.PrintErr("プロファイル取得エラー", err)
		return
	}

	// プロファイルを更新する
	profile2.SysGame = system_game_id2
	profile2.AdmGame = admin_game_id

	// プロファイルを保存する
	err = SaveProfile(profile2)

	// エラー処理
	if err != nil {
		logger.PrintErr("プロファイル保存エラー", err)
		return
	}

	logger.Println("プロファイル保存成功")
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

