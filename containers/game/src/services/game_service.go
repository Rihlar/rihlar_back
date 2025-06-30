package services

import (
	"game/logger"
	"game/models"

	"github.com/google/uuid"
)

type GameService struct{}

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