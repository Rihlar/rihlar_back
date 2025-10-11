package models

import (
	"game/logger"
)

// テーブル定義
type ItemBox struct {
	UserID   string `gorm:"primaryKey" json:"userID"`  // ユーザーID
	ItemID   string `gorm:"primaryKey" json:"itemID"`  // アイテムID
	Quantity int    `gorm:"default:0" json:"quantity"` // 所持数
}

// テーブル名
func (ItemBox) TableName() string {
	return "item_boxes"
}

func DebugItemBox() {

	userid := "user-e3abf90d-4bcf-4c3b-bbde-37694b1611b3"
	itemid := "item-e4c0e536-5cbe-4845-ac5d-3cd82dd86a15"

	// 書き込み
	result := Dbconn.Save(&ItemBox{
		UserID:   userid,
		ItemID:   itemid,
		Quantity: 5,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("アイテムボックス保存エラー", result.Error)
		return
	}

	logger.Println("アイテムボックス保存成功")

	// 取得コード
	returnData := ItemBox{}

	// 取得
	result = Dbconn.Where(&ItemBox{
		UserID: userid,
		ItemID: itemid,
	}).First(&returnData)

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("アイテムボックス取得エラー", result.Error)
		return
	}

	logger.Println("アイテムボックス取得成功")
}

// アイテムボックス取得
func GetItemBox(userId string) ([]ItemBox, error) {
	var itemBox []ItemBox

	result := Dbconn.Where("user_id = ?", userId).Find(&itemBox)
	if result.Error != nil {
		logger.PrintErr("アイテムボックス取得エラー", result.Error)
		return nil, result.Error
	}

	return itemBox, nil
}

// 所持しているか取得
func GetItemBoxByUserAndItem(userID, itemID string) (*ItemBox, error) {
    var boxes []ItemBox
    result := Dbconn.Where("user_id = ? AND item_id = ?", userID, itemID).Find(&boxes)
    if result.Error != nil {
        return nil, result.Error
    }

    if len(boxes) == 0 {
        return nil, nil // 存在しなければ nil
    }

    return &boxes[0], nil // 1件目を返す
}

// 新規作成
func CreateItemBox(box *ItemBox) error {
	return Dbconn.Create(box).Error
}

// 更新
func UpdateItemBox(box *ItemBox) error {
	return Dbconn.Save(box).Error
}
