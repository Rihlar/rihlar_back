package services

import (
	"gcore/logger"
	"gcore/models"
)

func Debug() {
	// CreateTestData()
	// メンバーを取得する
	// member, err := models.GetGame(models.AdminUserId1)
}

func CreateTestData() {
	for _, history := range models.AllUsersActionHistory {
		ReportMovement(MovementArgs{
			UserID:    history.UserID,
			Steps:     int64(history.Steps),
			Latitude:  history.Latitude,
			Longitude: history.Longitude,
			TimeStamp: history.Timestamp.Unix(),
		})
	}

	// admin ゲームを取得する
	admGame, err := models.GetGame(models.AdminGameId1)

	// エラー処理
	if err != nil {
		logger.PrintErr("ゲーム取得エラー", err)
		return
	}

	for _, circle := range models.CircleDatas {
		logger.Println("歩いた歩数", circle.Steps)

		ProcessCreateCircle(GamesCreateCircleArgs{
			UserID:     circle.UserID,
			Steps:      int64(circle.Steps),
			Latitude:   circle.Latitude,
			Longitude:  circle.Longitude,
			Games:      []models.Game{admGame},
			CircleSize: StepToSize(int64(circle.Steps)),
		})
	}
}


