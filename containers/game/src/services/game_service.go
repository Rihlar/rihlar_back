package services

import (
	"game/logger"
	"game/models"

	"github.com/google/uuid"
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

	// ゲームの情報を取得
	games, err := models.GetGame(allGames)
		if err != nil {
		logger.PrintErr("Game does not exist", err)
		return []models.Game{}, err
	}

	return games, nil
}

// adminゲーム作成
func(GameService) CreateAdminGame(bGame models.Game) error {

	// game_idを生成
	gameId, err := uuid.NewRandom() //新しいuuidの作成
	if err != nil {
		return err
	}

	// game_idのセット
	bGame.GameID = gameId.String()

	// 登録処理
	_, err = models.CreateAdminGame(&bGame)
		if err != nil {
		logger.PrintErr("Creation failure", err)
		return err
	}

	return err
}