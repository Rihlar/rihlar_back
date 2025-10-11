package controllers

import (
	"net/http"
	"game/models"

	"github.com/labstack/echo/v4"
)

// ユーザーが参加しているゲームの一覧を取得する
func GetUserParticipatedGamesHandler(c echo.Context) error {
	// URLからユーザーIDを取得
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "User ID is required"})
	}

	// モデルからゲーム一覧を取得
	games, err := models.GetGamesByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get user games"})
	}

	return c.JSON(http.StatusOK, games)
}

// 特定のゲームの行動ログを取得する
func GetMovementLogsHandler(c echo.Context) error {
	// URLからゲームIDとユーザーIDを取得
	gameID := c.Param("game_id")
	userID := c.Param("user_id")

	if gameID == "" || userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Game ID and User ID are required"})
	}

	// モデルから行動ログを取得
	logs, err := models.GetMovementLogs(gameID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get movement logs"})
	}

	return c.JSON(http.StatusOK, logs)
}
