package models

import "gcore/logger"

// テーブル定義
type GameChunk struct {
	ChunkID string `gorm:"primaryKey" json:"chunkID"`  // チャンクID
	GameID  string `gorm:"varchar(36)" json:"gameID"`  // ゲームID
	ImageID string `gorm:"varchar(36)" json:"imageID"` // イメージID
	OwnerID string `gorm:"varchar(36)" json:"ownerID"` // オーナーID
	Level   int    `gorm:"not null" json:"level"`      // 防衛レベル
}

// テーブル名
func (GameChunk) TableName() string {
	return "game_chunks"
}

// ゲームチャンク取得
func GetGameChunk(chunkid string) (GameChunk, error) {
	// 取得コード
	returnData := GameChunk{}

	// 取得する
	result := dbconn.Where(&GameChunk{
		ChunkID: chunkid,
	}).First(&returnData)

	return returnData, result.Error
}

// 所有者を変更する
func (gc *GameChunk) ChangeOwner(ownerid string) error {
	// 変更
	gc.OwnerID = ownerid

	// 更新
	return dbconn.Model(gc).Update("owner_id", ownerid).Error
}

// レベルを変更する
func (gc *GameChunk) ChangeLevel(level int) error {
	// 変更
	gc.Level = level

	// 更新
	return dbconn.Model(gc).Update("level", level).Error
}

func DebugGameChunk() {
	// デバッグ用のコードをここに書く

	chunkid := "3325d4ee-ef32-42a3-91d1-33d3582dffc2"
	gameid := "f36eb7ce-4e24-4805-99a5-b3ae3468708a"
	imageid := "76bd1e16-3105-4916-ad6b-7da9554c9601"
	ownerid := "e9178c88-3b64-4e61-b823-fd874d177d3c"

	// 書き込み
	result := dbconn.Save(&GameChunk{
		ChunkID: chunkid,
		GameID:  gameid,
		ImageID: imageid,
		OwnerID: ownerid,
		Level:   1,
	})

	// エラー処理
	if result.Error != nil {
		logger.PrintErr("ゲームチャンク保存エラー", result.Error)
		return
	}

	logger.Println("ゲームチャンク保存成功")

	// 取得する
	returnData, err := GetGameChunk(chunkid)

	// エラー処理
	if err != nil {
		logger.PrintErr("ゲームチャンク取得エラー", result)
		return
	}

	// 表示
	logger.Println("ゲームチャンク情報")
	logger.Println("チャンクID:", returnData.ChunkID)
	logger.Println("ゲームID:", returnData.GameID)
	logger.Println("イメージID:", returnData.ImageID)
	logger.Println("オーナーID:", returnData.OwnerID)
	logger.Println("レベル:", returnData.Level)

	// レベル変更
	err = returnData.ChangeLevel(2)
	if err != nil {
		logger.PrintErr("レベル変更エラー", err)
		return
	}

	//もう一度取得
	returnData, err = GetGameChunk(chunkid)
	if err != nil {
		logger.PrintErr("ゲームチャンク取得エラー", result)
		return
	}

	// 表示
	logger.Println("レベル変更後")
	logger.Println("レベル:", returnData.Level)

	// 所有者変更
	err = returnData.ChangeOwner("e9178c88-3b64-4e61-b823-fd874d177d3c")
	if err != nil {
		logger.PrintErr("所有者変更エラー", err)
		return
	}

	//もう一度取得
	returnData, err = GetGameChunk(chunkid)
	if err != nil {
		logger.PrintErr("ゲームチャンク取得エラー", result)
		return
	}

	// 表示
	logger.Println("所有者変更後")
	logger.Println("オーナーID:", returnData.OwnerID)

	logger.Println("ゲームチャンク取得成功")
}


