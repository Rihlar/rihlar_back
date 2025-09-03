package services

import (
	"fmt"
	"game/logger"
	"game/models"
	"math/rand"
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

// アイテムガチャ
func (ItemService) GetItemGacha(userId string) (models.Item, error) {
	// ガチャコスト
	cost := 100

	// 所持コイン確認
	coin, err := models.GetUserCoins(userId)
	if err != nil {
		return models.Item{}, err
	}

	if coin < cost {
		logger.PrintErr("コイン不足", err)
		return models.Item{}, fmt.Errorf("コイン不足でガチャを回せません")
	}

	// コインを減らす
	newAmount := coin - cost
	if err := models.UpdateUserCoins(userId, newAmount); err != nil {
		return models.Item{}, err
	}

	var items []models.Item
	// 全アイテム取得
	items, err = models.GetAllItems()
	if err != nil || len(items) == 0 {
		logger.PrintErr("items does not exist", err)
		return models.Item{}, err
	}

	// 全アイテムからランダムで1つ選ぶ
	idx := rand.Intn(len(items))
	chosen := items[idx]

	// 既存レコードを取得
	box, err := models.GetItemBoxByUserAndItem(userId, chosen.ItemID)
	if err != nil {
		return models.Item{}, err
	}

	if box == nil {
		// レコードがなければ新規作成
		box = &models.ItemBox{
			UserID:   userId,
			ItemID:   chosen.ItemID,
			Quantity: 1,
		}
		if err := models.CreateItemBox(box); err != nil {
			return models.Item{}, err
		}
	} else {
		// 既存レコードならQuantityを増やす
		box.Quantity += 1
		if err := models.UpdateItemBox(box); err != nil {
			return models.Item{}, err
		}
	}

	return chosen, nil
}
