package controllers

import (
	"game/logger"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUserParticipatedGamesHandler はユーザーが参加したゲーム一覧を取得するハンドラです。
func GetUserParticipatedGamesHandler(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id is required"})
	}

	gameIDs, err := services.GetUserParticipatedGames(userID)
	if err != nil {
		logger.PrintErr("GetUserParticipatedGamesHandler: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get participated games"})
	}

	return c.JSON(http.StatusOK, echo.Map{"gameIDs": gameIDs})
}

// GetMovementLogsHandler は特定のゲームの行動ログを取得するハンドラです。
func GetMovementLogsHandler(c echo.Context) error {
	userID := c.Param("user_id")
	gameID := c.Param("game_id")

	if userID == "" || gameID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "user_id and game_id are required"})
	}

	movementLogs, err := services.GetMovementLogs(userID, gameID)
	if err != nil {
		logger.PrintErr("GetMovementLogsHandler: ", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get movement logs"})
	}

	return c.JSON(http.StatusOK, echo.Map{"movementLogs": movementLogs})
}
