package services

import (
	"errors"
	"game/logger"
	"game/models"
	"game/utils"
	"time"
)

type GameService struct{}

// 終了済みゲーム一覧
func (GameService) GetEndGames(userId string) ([]models.Game, error) {

	// 全てのゲームを取得してくる
	allGames, err := models.GetJoinGames(userId)
	if err != nil {
		logger.PrintErr("Game ID does not exist", err)
		return []models.Game{}, err
	}

	// 終了済みゲームの一覧取得
	games, err := models.GetEndGames(allGames)
	if err != nil {
		logger.PrintErr("Unable to get game", err)
		return []models.Game{}, err
	}

	return games, nil
}

// 参加ゲーム一覧
func (GameService) GetJoinGames(userId string) ([]models.Game, error) {

	// 全てのゲームidを取得してくる(参加しているゲーム一覧取得)
	allGames, err := models.GetJoinGames(userId)
	if err != nil {
		logger.PrintErr("Game ID does not exist", err)
		return []models.Game{}, err
	}

	logger.Println(allGames)

	// 現在開催中ゲームの情報を取得
	games, err := models.GetGameHolding(allGames)
	if err != nil {
		logger.PrintErr("Game does not exist", err)
		return []models.Game{}, err
	}

	return games, nil
}

// ゲームに参加するエンドポイント
func (GameService) JoinGame(userId string, gameId string) error {
	// メンバーが存在するか判定
	exists, err := models.ExistsMember(userId, gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 存在している場合 (参加済みの場合)
	if exists {
		// エラーを返す
		return errors.New("you have already joined this game")
	}

	// プロファイル取得
	profile, err := models.GetProfile(userId)

	// エラー処理
	if err != nil {
		return err
	}

	// ゲームを更新
	profile.AdmGame = gameId

	// 更新
	err = models.SaveProfile(profile)

	// エラー処理
	if err != nil {
		return err
	}

	// チームのIDを生成する
	teamId, _ := utils.Genid()

	team := models.Team{
		TeamID: "teamid-" + teamId,
		GameID: gameId,
		Points: 0,
	}

	// チームを作成する
	err = models.CreateTeam(team)

	// エラー処理
	if err != nil {
		return err
	}

	// ゲームを取得する
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// チームを追加する
	err = game.AddTeam(team)

	// エラー処理
	if err != nil {
		return err
	}

	// チームにメンバーを追加
	err = team.AddMember(models.Member{
		GameID: gameId,
		TeamID: teamId,
		UserID: userId,
		Points: 0,
	})

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

type CreateGameArgs struct {
	Name         string `json:"name"`
	RegionID     string `json:"region_id"`
	StartTime    int64  `json:"start_time"`
	DulationDate int    `json:"dulation_date"`
}

func (GameService) CreateGame(args CreateGameArgs) error {
	// リージョンを取得する
	region, err := models.GetRegionByID(args.RegionID)

	// エラー処理
	if err != nil {
		return err
	}

	// ゲームIDを生成する
	gameid, err := utils.Genid()
	if err != nil {
		return err
	}

	// リージョンIDを更新する
	args.RegionID = region.RegionID

	// unix 時間を変換する
	startTime := time.Unix(0, args.StartTime)
	endTime := startTime.AddDate(0, 0, args.DulationDate)

	// ゲームを作成する
	return models.CreateGame(models.Game{
		GameID:    "gameid-" + gameid,
		StartTime: startTime,
		EndTime:   endTime,
		Flag:      0,
		Type:      0,
		Status:    0,
		RegionID:  args.RegionID,
	})
}

type GameMember struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Points   int    `json:"points"`
}

type GameTeam struct {
	TeamID string `json:"team_id"`
	Points int    `json:"points"`
}

type GameData struct {
	GameID       string       `json:"game_id"`
	Name         string       `json:"name"`
	RegionID     string       `json:"region_id"`
	Status       int          `json:"status"`
	StartTime    int64        `json:"start_time"`
	DulationDate int          `json:"dulation_date"`
	Members      []GameMember `json:"members"`
	Teams        []GameTeam   `json:"teams"`
}

// ゲームのリストを取得
func (GameService) GetGameList() ([]GameData, error) {
	// 全てのゲームを取得
	games,err := models.GetAllGames()

	// エラー処理
	if err != nil {
		return []GameData{}, err
	}

	returnData := []GameData{}

	for _, game := range games {
		// ゲームがシステムゲームの場合無視
		if game.Type == 0 {
			continue
		}

		// 日数の差分を取得
		diff := game.EndTime.Sub(game.StartTime)

		// チームを取得
		teams, err := getTeamFromGame(game)

		// エラー処理
		if err != nil {
			return []GameData{}, err
		}

		// メンバーを取得
		members, err := getMembersFromGame(game)

		// エラー処理
		if err != nil {
			return []GameData{}, err
		}

		returnData = append(returnData, GameData{
			GameID:       game.GameID,
			Name:         game.GameID,
			RegionID:     game.RegionID,
			StartTime:    game.StartTime.Unix(),
			Status:       game.Status,
			DulationDate: int(diff.Hours()) / 24,
			Members:      members,
			Teams:        teams,
		})
	}

	return returnData, nil
}

// ゲームのチームを取得
func getTeamFromGame(game models.Game) ([]GameTeam, error) {
	// チームを取得
	teams, err := game.GetTeams()

	// エラー処理
	if err != nil {
		return []GameTeam{}, err
	}

	returnTeams := []GameTeam{}

	for _, team := range teams {
		returnTeams = append(returnTeams, GameTeam{
			TeamID: team.TeamID,
			Points: team.Points,
		})
	}

	return returnTeams, nil
}

// メンバーを取得
func getMembersFromGame(game models.Game) ([]GameMember, error) {
	// チームを取得
	members, err := game.GetMembers()

	// エラー処理
	if err != nil {
		return []GameMember{}, err
	}

	returnMembers := []GameMember{}

	for _, member := range members {
		returnMembers = append(returnMembers, GameMember{
			UserID:   member.UserID,
			UserName: "",
			Points:   member.Points,
		})
	}

	return returnMembers, nil
}

// ゲームを削除
func (GameService) DeleteGame(gameId string) error {
	// ゲームを取得
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 削除する
	return game.DeleteGame()
}

// ゲームを開始
func (GameService) StartGame(gameId string) (error) {
	// ゲームを取得
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 開始する
	return game.StartGame()
}

// ゲームを終了
func (GameService) EndGame(gameId string) (error) {
	// ゲームを取得
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 終了する
	return game.EndGame()
}

// チームを削除するエンドポイント
func (GameService) DeleteTeam(gameId string, teamId string) (error) {
	// ゲームを取得
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 削除する
	return game.DeleteTeam(teamId)
}

// メンバーを削除するエンドポイント
func (GameService) DeleteMember(gameId string, userId string) (error) {
	// ゲームを取得
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 削除する
	return game.DeleteMember(userId)
}