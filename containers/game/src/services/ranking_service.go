package services

import (
	"game/logger"
	"game/models"
)

type RnakingService struct{}

// ランキングを取得し返還する用のテーブル
type RankingResult struct {
	UserID   string `json:"userID"`                 // ユーザーID
	Points   int    `gorm:"not null" json:"points"` // ポイント
	Rankings int    `json:"ranking"`                //ランキング
}

// 自分のランキング取得
func (RnakingService) GetMyRanking(userId string) (RankingResult, error) {

	// 全てのゲームを取得してくる
	allGames, err := models.GetPlaingGames(userId)
	if err != nil {
		logger.PrintErr("Game ID does not exist", err)
		return RankingResult{}, err
	}

		logger.Println(allGames)

	// 開催中ゲームの一覧取得
	games, err := models.GetGameHolding(allGames)
	if err != nil {
		logger.PrintErr("Unable to get game", err)
		return RankingResult{}, err
	}

	logger.Println(games)

	var gameId string
	// adminゲームか判断してIDを保持する
	for _, game := range games {
		logger.Println("あああああ", game)
		// adminゲームはTypeが１
		if game.Type == 1 {
			gameId = game.GameID
			logger.Println("あああああ")
		}
	}

	logger.Println(gameId) 

	// 特定したgameIdとuserIdからランキングを取得
	ranking, err := models.GetMyRanking(userId, gameId)
	if err != nil {
		logger.PrintErr("Unable to get ranking", err)
		return RankingResult{}, err
	}
		logger.Println(ranking)

	// 自己満で得点も返したいので取得してくる
	user, err := models.GetMyPoints(userId, gameId)
	if err != nil {
		logger.PrintErr("Can't get points", err)
		return RankingResult{}, err
	}

	// 返す用のデータ
	result := RankingResult{
		UserID:   userId,
		Points:   user.Points,
		Rankings: ranking,
	}

	return result, nil
}
