package models

import "game/logger"

// テーブル定義
type Item struct {
	ItemID      string `gorm:"primaryKey" json:"itemID"`       // アイテムID
	ItemName    string `gorm:"varchar(50)" json:"itemName"`    // アイテム名
	Description string `gorm:"varchar(50)" json:"description"` // アイテム説明
}

// テーブル名
func (Item) TableName() string {
	return "items"
}

// テストデータの宣言
var (
	// TODO とりあえず現状のリージョンをかく
	items = []Item{
		{
			ItemID:      "item-e4c0e536-5cbe-4845-ac5d-3cd82dd86a15",
			ItemName:    "多分守るマン",
			Description: "円を守る(効果時間は2時間)",
		},
		{
			ItemID:      "item-494059a1-dc63-41bb-9b6b-29d761b6ae3b",
			ItemName:    "倍倍ふぁいと",
			Description: "円のポイント2倍",
		},
		{
			ItemID:      "item-117bbf40-14c0-4182-b120-124770c732f0",
			ItemName:    "嬉Cー",
			Description: "コインの獲得量2倍",
		},
		{
			ItemID:      "item-c60ac786-68eb-4616-8b64-2fb56f7a4f3b",
			ItemName:    "絶対押してマイルマン",
			Description: "円を絶対に取れる多分守るマンより強い",
		},
		{
			ItemID:      "item-26164e99-790e-429f-95df-8008eb34d628",
			ItemName:    "一向に構わん水",
			Description: "全ての状態異常をキャンセル自分に向けたバフも消滅",
		},
		{
			ItemID:      "item-aea0e699-ec54-42b3-ad7a-f34e32c45d18",
			ItemName:    "勘のいいガキー",
			Description: "運営からの指定写真のヒントが他の人より先に見える",
		},
		{
			ItemID:      "item-c4725712-f7cf-4b27-8f37-9c806f4deada",
			ItemName:    "命を刈り取るカマ",
			Description: "相手のポイント削る",
		},
	}
)


func DebugItem() {
	// アイテムの登録する
	for _, item := range items {
		// 書き込み
		err := Dbconn.Save(&item).Error
		// エラー処理
		if err != nil {
			logger.PrintErr("アイテム登録エラー", err)
			return
		}
	}

	logger.Println("アイテム登録成功")
}

func GetItemDeteile(itemId string) (Item, error) {
	var item Item

	result := Dbconn.Where("item_id = ?", itemId).Take(&item)
	if result.Error != nil {
		logger.PrintErr("アイテム詳細取得エラー", result.Error)
		return Item{}, result.Error
	}

	return item, nil
}

func GetAllItems() ([]Item, error) {
	var items []Item
    // 全アイテム取得
    err := Dbconn.Find(&items).Error
	if  err != nil {
        return []Item{}, err
    }

	return items, nil
}