package services

import (
	"game/logger"
	"game/models"
)

type CircleService struct{}

func(CircleService) GetCircleDeteile(circleId string) (models.Circle, error) {

	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
		return models.Circle{}, err
	}

	return circleDeteile, nil
}