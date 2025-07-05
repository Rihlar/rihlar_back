package services

import (
	"errors"
	"game/logger"
	"game/models"
	"game/utils"
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
	teamId,_ := utils.Genid()
	
	team := models.Team{
		TeamID:    "teamid-" + teamId,
		GameID:    gameId,
		Points:    0,
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
