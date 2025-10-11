package services

import (
	"game/logger"
	"game/models"

	"gorm.io/gorm"
)

// GetUserParticipatedGames はユーザーが参加したゲーム一覧を取得します。
func GetUserParticipatedGames(userID string) ([]string, error) {
	var gameIDs []string
	err := models.dbconn.Model(&models.MovementLog{}). // models.dbconn を使用
		Where("user_id = ?", userID).
		Distinct("game_id").
		Pluck("game_id", &gameIDs).Error
	if err != nil {
		logger.PrintErr("ユーザーが参加したゲーム一覧の取得に失敗しました", err)
		return nil, err
	}
	return gameIDs, nil
}

// GetMovementLogs は特定のゲームの行動ログを取得します。
func GetMovementLogs(userID, gameID string) ([]models.MovementLog, error) {
	var movementLogs []models.MovementLog
	err := models.dbconn.Where("user_id = ? AND game_id = ?", userID, gameID).
		Order("time_stamp asc").
		Find(&movementLogs).Error
	if err != nil {
		logger.PrintErr("行動ログの取得に失敗しました", err)
		return nil, err
	}
	return movementLogs, nil
}
