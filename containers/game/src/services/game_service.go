package services

import (
	"errors"
	"game/logger"
	"game/models"
	"game/utils"
	"net/http"
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
		Type:      1,
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
	games, err := models.GetAllGames()

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

type AllGameData struct {
	IsJoined  bool   `json:"isJoined"`
	GameID    string `json:"gameID"`
	StartTime int64 `json:"startTime"`
	EndTime   int64 `json:"endTime"`
	Flag      int    `json:"flag"`
	Type      int    `json:"type"`
	Status    int    `json:"status"`
	RegionID  string `json:"regionID"`
}

// すべてのゲームのリストを取得
func (GameService) GetAllGameList(userid string) ([]AllGameData, error) {
	// 全てのゲームを取得
	games, err := models.GetAllGames()

	// エラー処理
	if err != nil {
		return []AllGameData{}, err
	}

	returnData := []AllGameData{}

	for _, game := range games {
		// ゲームがシステムゲームの場合無視
		if game.Type == 0 {
			continue
		}

		returnData = append(returnData, AllGameData{
			IsJoined:  game.CheckJoin(userid),
			GameID:    game.GameID,
			StartTime: game.StartTime.Unix(),
			EndTime:   game.EndTime.Unix(),
			Flag:      game.Flag,
			Type:      game.Type,
			Status:    game.Status,
			RegionID:  game.RegionID,
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
func (GameService) StartGame(gameId string) error {
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
func (GameService) EndGame(gameId string) error {
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
func (GameService) DeleteTeam(gameId string, teamId string) error {
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
func (GameService) DeleteMember(gameId string, userId string) error {
	// ゲームを取得
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return err
	}

	// 削除する
	return game.DeleteMember(userId)
}

type StartedGameResult struct {
	GameID    string `json:"game_id"`    // ゲームID
	StartedAt int64  `json:"started_at"` // 開始日
	RegionID  string `json:"region_id"`  // 開催地
	EndAt     int64  `json:"end_at"`     // 終了日
	Players   int    `json:"players"`    // 参加人数
	IsJoined  bool   `json:"is_joined"`  // ゲームに参加しているか
}

// ユーザーのリージョンから開催中のゲームを取得する完遂
func (GameService) GetStartedGames(userId string) (StartedGameResult, CustomError) {
	// プロフィールを取得する
	profile, err := models.GetProfile(userId)
	if err != nil {
		return StartedGameResult{}, CustomError{
			Code:       http.StatusInternalServerError,
			LogMessage: "プロフィール取得エラー",
			ErrMessage: "プロフィール取得エラー",
			Err:        err,
		}
	}

	logger.Println("profile", profile)

	// ゲームを検索する
	games, err := models.SearchGame(models.SearchGameArgs{
		IsSearchSystem: false,
		IsSearchRegion: true,
		RegionID:       profile.RegionID,
		IsSearchStatus: true,
		Status:         1,
	})

	// エラー処理
	if err != nil {
		return StartedGameResult{}, CustomError{
			Code:       http.StatusInternalServerError,
			LogMessage: "ゲーム検索エラー",
			ErrMessage: "ゲーム検索エラー",
			Err:        err,
		}
	}

	// ゲームがない場合
	if len(games) == 0 {
		return StartedGameResult{
				IsJoined: false,
			}, CustomError{
				Code:       http.StatusNotFound,
				LogMessage: "開催中のゲームがありません",
				ErrMessage: "開催中のゲームがありません",
				Err:        errors.New("開催中のゲームがありません"),
			}
	}

	// 現状一つだけなので
	game := games[0]

	// 全てのメンバー取得
	members, err := game.GetMembers()

	// エラー処理
	if err != nil {
		return StartedGameResult{}, CustomError{
			Code:       http.StatusInternalServerError,
			LogMessage: "メンバー取得エラー",
			ErrMessage: "メンバー取得エラー",
			Err:        err,
		}
	}

	// メンバーの人数を計算
	playerCount := len(members)

	return StartedGameResult{
		GameID:    game.GameID,
		StartedAt: game.StartTime.Unix(),
		RegionID:  game.RegionID,
		EndAt:     game.EndTime.Unix(),
		Players:   playerCount,
		IsJoined:  true,
	}, CustomError{}
}

// GameData represents the overall structure of the game data.
type MySelfGameData struct {
	IsAdminJoined bool       `json:"IsAdminJoined"`
	Admin         AdminData  `json:"admin"`
	System        SystemData `json:"system"`
}

// AdminData holds administrative information about the game.
type AdminData struct {
	IsFinished bool   `json:"IsFinished"` // ゲームが終了済みか
	IsStarted  bool   `json:"IsStarted"`  // ゲームが開始済みか
	GameID     string `json:"GameID"`     // ゲームID
	StartTime  int64  `json:"StartTime"`  // 開始時間 Unix Time
	EndTime    int64  `json:"EndTime"`    // 終了時間 Unix Time
}

// SystemData holds system-level information about the game.
type SystemData struct {
	GameID string `json:"GameID"` // ゲームID
}

// 自身が参加しているゲームを取得する
func (GameService) GetMyGames(userId string) (MySelfGameData, error) {
	// メンバー一覧を取得する
	gameIds, err := models.GetJoinGames(userId)
	if err != nil {
		return MySelfGameData{}, err
	}

	returnData := MySelfGameData{}

	// ゲームIDを回す
	for _, gameId := range gameIds {
		// ゲームを取得する
		game, err := models.GetGame(gameId)
		if err != nil {
			logger.Println(err)
			continue
		}

		if game.Type == 1 {
			isFinished := false

			if game.Status == 2 {
				// 終了の場合
				isFinished = true
			}

			isStarted := false

			if game.Status == 1 {
				// 開始の場合
				isStarted = true
			}

			// admin の場合返すデータに
			returnData.Admin = AdminData{
				IsFinished: isFinished,
				IsStarted:  isStarted,
				GameID:     gameId,
				StartTime:  game.StartTime.Unix(),
				EndTime:    game.EndTime.Unix(),
			}

			returnData.IsAdminJoined = true
		} else {
			// system の場合返すデータに
			returnData.System = SystemData{
				GameID: gameId,
			}
		}
	}

	return returnData, nil
}
