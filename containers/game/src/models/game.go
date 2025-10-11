package models

import (
	"game/logger"
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

// デバック用
func DebugGame() {
	// ゲームを作成する関数
	CreateGame(Game{
		GameID:    "gameid-996e5916-28b7-4222-ad5c-b332c1f892ec",
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0,0,20),
		Flag:      0,
		Type:      1,
		Status:    1,
		RegionID:  "regionId-ef5aa179-53e0-481d-b64d-ae7654049a88",
	})
}

// ゲームを作成する関数
func CreateGame(game Game) error {
	// 書き込み
	return Dbconn.Save(&game).Error
}

// TODO デバッグ用 チームを追加する
func (game *Game) AddTeam(team Team) error {
	return Dbconn.Model(game).Association("Teams").Append(&team)
}

// ゲームを消す関数
func (game *Game) DeleteGame() error {
	return Dbconn.Delete(game).Error
}

// 全てのゲームを取得
func GetAllGames() ([]Game, error) {
	// 結果格納用
	var games []Game

	// 取得する
	err := Dbconn.Find(&games).Error

	// エラー処理
	if err != nil {
		return []Game{}, err
	}

	return games, nil
}

// ユーザーがゲームに参加しているかを判定
func (game *Game) CheckJoin(userId string) bool {
	// メンバーを検索する
	_, err := game.GetMemberByUserID(userId)

	// エラー処理
	if err != nil {
		return false
	}

	return true
}

type SearchGameArgs struct {
	IsSearchSystem bool	// システムかどうか

	IsSearchRegion bool	// リージョン検索
	RegionID string

	IsSearchStatus bool	// ステータス検索
	Status   int
}

// ゲームを検索する関数
func SearchGame(args SearchGameArgs) ([]Game, error) {
	// 検索用
	searchParam := Game{}

	// 検索条件を設定する
	if args.IsSearchRegion {
		searchParam.RegionID = args.RegionID
	}

	// 検索条件を設定する
	if args.IsSearchStatus {
		searchParam.Status = args.Status
	}

	if args.IsSearchSystem {
		// システムゲームを検索
		searchParam.Type = 0
	} else {
		// 管理者ゲームを検索
		searchParam.Type = 1
	}

	// 結果格納用
	var games []Game

	// 取得する
	err := Dbconn.Debug().Where(&searchParam).Find(&games).Error

	// エラー処理
	if err != nil {
		return []Game{}, err
	}

	return games, nil
}

// ゲームを開始する
func (game *Game) StartGame() error {
	return Dbconn.Debug().Model(game).Update("status", 1).Error
}

// ゲームを終了する
func (game *Game) EndGame() error {
	return Dbconn.Debug().Model(game).Update("status", 2).Error
}

// ゲームに属する全てのメンバー取得
func (game *Game) GetMembers() ([]Member, error) {
	// 結果格納用
	var members []Member

	// 取得する
	err := Dbconn.Where(&Member{
		GameID: game.GameID,
	}).Find(&members).Error

	// エラー処理
	if err != nil {
		return []Member{}, err
	}

	return members, nil
}

func (game *Game) GetTeams() ([]Team, error) {
	// 結果格納用
	var teams []Team

	// 取得する
	err := Dbconn.Where(&Team{
		GameID: game.GameID,
	}).Find(&teams).Error

	// エラー処理
	if err != nil {
		return []Team{}, err
	}

	return teams, nil
}

// メンバーを取得
func (game *Game) GetMemberByUserID(userid string) (Member, error) {
	// 取得する
	returnData := Member{}

	// 取得する
	err := Dbconn.Where(&Member{
		UserID: userid,
		GameID: game.GameID,
	}).First(&returnData).Error

	// エラー処理
	if err != nil {
		return Member{}, err
	}

	return returnData, nil
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
	err = Dbconn.Where(&Team{
		TeamID: member.TeamID,
	}).First(&returnData).Error

	// エラー処理
	if err != nil {
		return Team{}, err
	}

	return returnData, nil
}

// ランキング上位取得
func (game *Game) GetRanking() ([]Team, error) {
	var rankings []Team

	result := Dbconn.
		Where("game_id = ?", game.GameID).
		Order("points DESC").
		Find(&rankings)

	if result.Error != nil {
		logger.PrintErr("ランキング上位取得エラー", result.Error)
		return nil, result.Error
	}

	return rankings, nil
}

// ゲームの詳細取得
func GetGames(gameId []string) ([]Game, error) {
	// 結果格納用
	var games []Game

	result := Dbconn.Where("game_id IN ?", gameId).Find(&games)
	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return []Game{}, nil
	}

	return games, nil
}

// ゲームを取得
func GetGame(gameId string) (Game, error) {
	var game Game

	// 取得する
	err := Dbconn.Where(&Game{
		GameID: gameId,
	}).First(&game).Error

	// エラー処理
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

// 開催中のゲーム取得
func GetGameHolding(gameId []string) ([]Game, error) {
	// 結果格納用
	var games []Game

	result := Dbconn.Where("game_id IN ?", gameId).Where("status = ?", 1).Find(&games)
	if result.Error != nil {
			return []Game{}, nil
	}

	return games, nil
}

// 終了済みゲーム一覧取得
func GetEndGames(gameId []string) ([]Game, error) {
	// 結果格納用
	var games []Game

	// statusが2で絞る
	result := Dbconn.Where("game_id IN ?", gameId).Where("status = ?", 2).Find(&games)
	if result.Error != nil {
		logger.PrintErr("ゲーム取得エラー", result.Error)
		return []Game{}, nil
	}

	return games, nil
}

// ID からゲームを取得
func GetGameByID(gameId string) (Game, error) {
	var game Game

	result := Dbconn.Where(&Game{
		GameID: gameId,
	}).First(&game)

	if result.Error != nil {
		logger.PrintErr("ゲームID取得エラー", result.Error)
		return Game{}, result.Error
	}

	return game, nil
}

// ユーザーIDから参加しているゲーム一覧を取得する
func GetGamesByUserID(userID string) ([]Game, error) {
	// ユーザーが参加しているゲームID一覧を取得
	gameIDs, err := GetJoinGames(userID)
	if err != nil {
		logger.PrintErr("参加ゲームIDの取得エラー", err)
		return nil, err
	}

	if len(gameIDs) == 0 {
		// 参加しているゲームがない場合は空のスライスを返す
		return []Game{}, nil
	}

	// ゲームID一覧からゲーム情報を取得
	games, err := GetGames(gameIDs)
	if err != nil {
		logger.PrintErr("ゲーム情報の取得エラー", err)
		return nil, err
	}

	return games, nil
}

// Top10のランキングを取得
func (game *Game) GetRankingTop10() ([]Team, error) {
	var teams []Team

	result := Dbconn.
		Where("game_id = ?", game.GameID).
		Order("points DESC").
		Limit(10).
		Find(&teams)

	if result.Error != nil {
		logger.PrintErr("ランキング上位取得エラー", result.Error)
		return nil, result.Error
	}

	return teams, nil
}