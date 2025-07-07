package services

import (
	"game/logger"
	"game/models"
)

type CircleService struct{}

//　円の詳細取得
func(CircleService) GetCircleDeteile(circleId string) (models.Circle, error) {

	//　円の詳細取得
	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
		return models.Circle{}, err
	}

	return circleDeteile, nil
}

// 画像取得
func(CircleService) GetCircleImage(circleId string) (string, error) {

	// circleIdから円のimageIdをとってくる
	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
		return "", err
	}

	// imageIdから画像パス生成　TODO:
	imagePath := "./assets/circle-images/" + circleDeteile.ImageID + ".png"

	return imagePath, nil
}
