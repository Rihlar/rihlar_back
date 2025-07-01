package services

type Ranking struct {
	TopData  []RankingData `json:"topData"`
	SelfData RankingData `json:"selfData"`
}

type RankingData struct {
	TeamID string `json:"teamID"`
	Points int    `json:"points"`
}

func GetRanking(gameID string, userID string) ([]Ranking, error) {
	

	return nil, nil
}