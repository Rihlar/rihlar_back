package services

import (
	"game/models"
)

type TeamService struct{}

func (s *TeamService) GetTeams(gameID string) ([]models.Team, error) {
	return models.GetTeamsByGameID(gameID)
}
