package services

import (
	"gcore/logger"
	"gcore/models"
)

type Ranking struct {
	TopData  []RankingData `json:"topData"`
	SelfData RankingData `json:"selfData"`
}

type RankingData struct {
	TeamID string `json:"teamID"`
	Points int    `json:"points"`
}

func GetRanking(gameID string, userID string) (Ranking, error) {
	// ゲームを取得する
	game, err := models.GetGame(gameID)

	// エラー処理
	if err != nil {
		return Ranking{}, err
	}

	// ランキングを取得する
	rankings, err := game.GetRanking(10)

	// エラー処理
	if err != nil {
		return Ranking{}, err
	}

	logger.Println("Rankings: ", rankings)

	returnData := Ranking{}

	for _, ranking := range rankings {
		// ランキングに追加する
		returnData.TopData = append(returnData.TopData, RankingData{
			TeamID: ranking.TeamID,
			Points: ranking.Points,
		})
	}

	// 自身を設定する
	// メンバー取得
	members, err := game.GetMemberByUserID(userID)

	// エラー処理
	if err != nil {
		return Ranking{}, err
	}

	// チームを取得する
	team, err := game.GetTeam(members.TeamID)

	// エラー処理
	if err != nil {
		return Ranking{}, err
	}

	// ランキングに追加する
	returnData.SelfData = RankingData{
		TeamID: team.TeamID,
		Points: team.Points,
	}

	return returnData, nil
}

