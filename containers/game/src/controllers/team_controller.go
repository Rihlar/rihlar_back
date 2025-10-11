package controllers

import (
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var teamService = services.TeamService{}

func GetTeamsHandler(c echo.Context) error {
	gameID := c.Param("game_id")
	if gameID == "" {
		return c.JSON(http.StatusBadRequest, "Game ID is required")
	}

	teams, err := teamService.GetTeams(gameID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, teams)
}
