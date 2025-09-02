package services

import (
	"game/logger"
	"game/models"
)

type ItemService struct{}


// アイテムボック取得
func (ItemService) GetItemBox(userId string) ([]models.ItemBox, error) {

	// アイテムボックス取得
	itemBox, err := models.GetItemBox(userId)
	if err != nil {
		logger.PrintErr("itembox does not exist", err)
		return []models.ItemBox{}, err
	}

	return itemBox, nil
}

// アイテム詳細取得
func (ItemService) GetItemDeteile(itemId string) (models.Item, error) {

	// アイテム詳細取得
	item, err := models.GetItemDeteile(itemId)
	if err != nil {
		logger.PrintErr("itembox does not exist", err)
		return models.Item{}, err
	}

	return item, nil
}