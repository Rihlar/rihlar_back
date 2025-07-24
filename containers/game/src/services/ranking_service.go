package services

import (
	"game/logger"
	"game/models"
)

type RankingService struct{}

// ランキングを取得し返還する用のテーブル
type RankingResult struct {
	UserID   string `json:"userID"`                 // ユーザーID
	Points   int    `gorm:"not null" json:"points"` // ポイント
	Rankings int    `json:"ranking"`                //ランキング
}

// 自分のランキング取得
func (RankingService) GetMyRanking(userId string) (RankingResult, error) {

	// 全てのゲームを取得してくる
	allGames, err := models.GetJoinGames(userId)
	if err != nil {
		logger.PrintErr("Game ID does not exist", err)
		return RankingResult{}, err
	}

	// 開催中ゲームの一覧取得
	games, err := models.GetGameHolding(allGames)
	if err != nil {
		logger.PrintErr("Unable to get game", err)
		return RankingResult{}, err
	}

	var gameId string
	// adminゲームか判断してIDを保持する
	for _, game := range games {
		// adminゲームはTypeが１
		if game.Type == 1 {
			gameId = game.GameID
		}
	}

	// 特定したgameIdとuserIdからランキングを取得
	ranking, err := models.GetMyRanking(userId, gameId)
	if err != nil {
		logger.PrintErr("Unable to get ranking", err)
		return RankingResult{}, err
	}

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

// TopRankingは全体のJSON構造を表します。
// JSONのトップレベルキーが固定されている場合に最適です。
type TopRanking struct {
	Top1  TopEntry `json:"Top1"`
	Top2  TopEntry `json:"Top2"`
	Top3  TopEntry `json:"Top3"`
	Other TopEntry `json:"Other"`
	Self  TopEntry `json:"Self"`
}

// TopEntryはTeamIDとCirclesのリストを持つエントリを表します。
type TopEntry struct {
	TeamID  string   `json:"TeamID"`
	Circles []Circle `json:"Circles"`
}

// Circleは個々の円のデータ（ID、位置情報、サイズ、レベルなど）を表します。
type Circle struct {
	CircleID  string  `json:"CircleID"`
	GameID    string  `json:"GameID"`
	Size      int     `json:"Size"`
	Level     int     `json:"Level"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	ImageID   string  `json:"ImageID"`
	TimeStamp int64   `json:"TimeStamp"` // Unixタイムスタンプ (秒)
}

// ランキングtop3と自信を取得
func (RankingService) GetRankingTop(userid, gameId string) (TopRanking, error) {
	// ゲームを取得
	game, err := models.GetGameByID(gameId)

	// エラー処理
	if err != nil {
		logger.PrintErr("Game does not exist", err)
		return TopRanking{}, err
	}

	// チームのランキングを取得
	rankings, err := game.GetRanking()
	if err != nil {
		logger.PrintErr("Unable to get ranking", err)
		return TopRanking{}, err
	}

	returnRanking := TopRanking{}

	// チーム一覧を回す
	for index, team := range rankings {
		// チームの円の一覧を取得する
		circles, err := models.GetCircleByTeamId(team.TeamID)

		// エラー処理
		if err != nil {
			logger.PrintErr("サークル取得エラー", err)
			continue
		}

		// 3位いないなら
		if index <= 3 {
			//TODO 力技なので今後修正する
			if index == 0 {
				returnRanking.Top1.TeamID = team.TeamID
				returnRanking.Top1.Circles = ModelCircleToCircle(circles, false)
			} else if index == 1 {
				returnRanking.Top2.TeamID = team.TeamID
				returnRanking.Top2.Circles = ModelCircleToCircle(circles, false)
			} else if index == 2 {
				returnRanking.Top3.TeamID = team.TeamID
				returnRanking.Top3.Circles = ModelCircleToCircle(circles, false)
			}

			// 戻る
			continue
		}

		// 3位以降はOtherに入れる
		returnRanking.Other.TeamID = ""
		returnRanking.Other.Circles = append(returnRanking.Other.Circles, ModelCircleToCircle(circles, false)...)
	}

	// チームを取得
	// メンバーを取得
	selfTeam, err := game.GetTeamByUserID(userid)

	// エラー処理
	if err != nil {
		return TopRanking{}, err
	}

	// チームの円の一覧を取得する
	circles, err := models.GetCircleByTeamId(selfTeam.TeamID)

	// エラー処理
	if err != nil {
		logger.PrintErr("サークル取得エラー", err)
		return TopRanking{}, err
	}

	// 円を取得して追加する
	returnRanking.Self.TeamID = selfTeam.TeamID
	returnRanking.Self.Circles = ModelCircleToCircle(circles, true)

	return returnRanking, nil
}

// ランキングTOP10用の構造体
type RankTop10Entry struct {
	TeamId string `json:"TeamID"`
	Points int    `json:"Points"`
}

type SelfTop10Entry struct {
	Rank   int    `json:"rank"`
	Point  int    `json:"point"`
	TeamId string `json:"TeamID"`
}

type Top10LeaderboardData struct { // 構造体名を変更
	Ranks []RankTop10Entry `json:"ranks"`
	Self  SelfTop10Entry   `json:"self"`
}

// ソロ用ランキングTOP10用の構造体
type SoloRankTop10Entry struct {
	TeamId   string `json:"TeamID"`
	UserName string `json:"UserName"`
	Points   int    `json:"Points"`
}

type SoloSelfTop10Entry struct {
	Rank     int    `json:"rank"`
	Point    int    `json:"point"`
	UserName string `json:"UserName"`
	TeamId   string `json:"TeamID"`
}

// ソロ用ランキングTOP10用の構造体
type SoloTop10LeaderboardData struct { // 構造体名を変更
	Ranks []SoloRankTop10Entry `json:"ranks"`
	Self  SoloSelfTop10Entry   `json:"self"`
}

// ランキングtop10を取得
func (RankingService) GetRankingTop10(userid, gameId string) (Top10LeaderboardData, error) {
	// ゲームを取得する
	game, err := models.GetGameByID(gameId)

	// エラー処理
	if err != nil {
		logger.PrintErr("Game does not exist", err)
		return Top10LeaderboardData{}, err
	}

	// メンバーを取得
	member, err := game.GetMemberByUserID(userid)

	// エラー処理
	if err != nil {
		return Top10LeaderboardData{}, err
	}

	// top10のランキングを取得
	teams, err := game.GetRankingTop10()
	if err != nil {
		logger.PrintErr("Unable to get ranking", err)
		return Top10LeaderboardData{}, err
	}

	datas := []RankTop10Entry{}

	// チームを回す
	for _, team := range teams {
		// チームを追加する
		datas = append(datas, RankTop10Entry{
			TeamId: team.TeamID,
			Points: team.Points,
		})
	}

	// 自身のチーム情報を取得する
	myTeam, err := member.GetTeam()

	// エラー処理
	if err != nil {
		return Top10LeaderboardData{}, err
	}

	// 自身のチームランキング取得
	myRank, err := myTeam.GetRank()
	if err != nil {
		return Top10LeaderboardData{}, err
	}

	return Top10LeaderboardData{
		Ranks: datas,
		Self: SelfTop10Entry{
			Rank:   myRank,
			Point:  myTeam.Points,
			TeamId: myTeam.TeamID,
		},
	}, nil
}

// ソロ用のランキングtop10を取得
func (RankingService) GetSoloRankingTop10(userid, gameId string) (SoloTop10LeaderboardData, error) {
	// ゲームを取得する
	game, err := models.GetGameByID(gameId)

	// エラー処理
	if err != nil {
		logger.PrintErr("Game does not exist", err)
		return SoloTop10LeaderboardData{}, err
	}

	// メンバーを取得
	member, err := game.GetMemberByUserID(userid)

	// エラー処理
	if err != nil {
		return SoloTop10LeaderboardData{}, err
	}

	// top10のランキングを取得
	teams, err := game.GetRankingTop10()
	if err != nil {
		logger.PrintErr("Unable to get ranking", err)
		return SoloTop10LeaderboardData{}, err
	}

	datas := []SoloRankTop10Entry{}

	// チームを回す
	for _, team := range teams {
		// メンバーを一人取得
		soloMemberId := team.Members[0].UserID

		// プロファイルを取得
		profile, err := models.GetProfile(soloMemberId)

		// 	エラー処理
		if err != nil {
			return SoloTop10LeaderboardData{}, err
		}

		// チームを追加する
		datas = append(datas, SoloRankTop10Entry{
			TeamId:   team.TeamID,
			Points:   team.Points,
			UserName: profile.Name,
		})
	}

	// 自身のチーム情報を取得する
	myTeam, err := member.GetTeam()

	// エラー処理
	if err != nil {
		return SoloTop10LeaderboardData{}, err
	}

	// 自身のチームランキング取得
	myRank, err := myTeam.GetRank()
	if err != nil {
		return SoloTop10LeaderboardData{}, err
	}

	// 自身のユーザー情報を取得する
	profile, err := models.GetProfile(userid)

	// エラー処理
	if err != nil {
		return SoloTop10LeaderboardData{}, err
	}

	return SoloTop10LeaderboardData{
		Ranks: datas,
		Self: SoloSelfTop10Entry{
			Rank:     myRank,
			Point:    myTeam.Points,
			TeamId:   myTeam.TeamID,
			UserName: profile.Name,
		},
	}, nil
}

// モデルの円を返す円に変換する
func ModelCircleToCircle(circles []models.Circle, isSelf bool) []Circle {
	returnCircles := []Circle{}

	for _, circle := range circles {
		// 自分以外の場合 ImageID を消す
		if !isSelf {
			circle.ImageID = ""
		}

		returnCircles = append(returnCircles, Circle{
			CircleID:  circle.CircleID,
			Size:      circle.Size,
			Level:     circle.Level,
			Latitude:  circle.Latitude,
			Longitude: circle.Longitude,
			ImageID:   circle.ImageID,
			TimeStamp: circle.CreatedAT.Unix(),
			GameID:    circle.GameID,
		})
	}

	return returnCircles
}
