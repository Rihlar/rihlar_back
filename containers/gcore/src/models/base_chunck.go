package models

import (
	"gcore/location"
	"gcore/logger"
	"gcore/utils"

	"gorm.io/gorm"
)

// テーブル定義
type BaseChunk struct {
	GridID    string  `gorm:"primaryKey" json:"gridID"`               // グリッドID
	StartLat  float64 `gorm:"not null" json:"startLat"`              // 開始緯度
	StartLon  float64 `gorm:"not null" json:"startLon"`              // 開始経度
	EndLat    float64 `gorm:"not null" json:"endLat"`                // 終了緯度
	EndLon    float64 `gorm:"not null" json:"endLon"`                // 終了経度
	RegionID  string  `gorm:"varchar(50);primaryKey" json:"regionID"` // ゲーム開催地域
}

// テーブル名
func (BaseChunk) TableName() string {
	return "chunks"
}

// TODO デバッグ用 ゲーム用のリージョンを作成する関数
func (game *Game) FillRegion(region Region) error {
	// グリッド生成
	grids := location.GenerateGrid(region.StartLat, region.StartLon, region.EndLat, region.EndLon, GridMeter)

	for _, grid := range grids {
		// チャンクをキャッシュに保存する
		err := location.SaveChunk(region.RegionID, grid.ID, grid.BottomRight.Lat, grid.BottomRight.Lon)

		// エラー処理
		if err != nil {
			return err
		}

		// チャンクを保存する
		err = dbconn.Create(&BaseChunk{
			StartLat:  grid.BottomRight.Lat,
			StartLon:  grid.BottomRight.Lon,
			EndLat:    grid.TopLeft.Lat,
			EndLon:    grid.TopLeft.Lon,
			RegionID:  region.RegionID,
			GridID:    grid.ID,
		}).Error

		// エラー処理
		if err != nil {
			logger.PrintErr("チャンク保存エラー", err)
			return err
		}
	}

	return nil
}

// 緯度経度からチャンクを取得する
func (game *Game) GetChunkByLatLon(lat, lon float64) (GameChunk, error) {
	// region の ID がキャッシュに保存されてる
	exists := location.ExistsRegion(game.RegionID)

	// 存在しない場合は埋める
	if !exists {
		// 埋める処理
		// リージョンを取得
		region, err := GetRegionByID(game.RegionID)
		if err != nil {
			return GameChunk{}, err
		}

		// リージョンを埋める
		err = game.FillRegion(region)
		if err != nil {
			return GameChunk{}, err
		}
	}

	// キャッシュから一番近いgrididを取得する
	gridId,err := location.FindNearChunk(game.RegionID, lat, lon, GridMeter * 2)
	if err != nil {
		return GameChunk{}, err
	}

	findChunk := GameChunk{}

	logger.Println("gridID", gridId)
	logger.Println("ゲームID", game.GameID)

	// ゲームからチャンクを取得する
	err = dbconn.Where(GameChunk{
		GameID:   game.GameID,
		GridID:   gridId,
	}).First(&findChunk).Error

	// 存在しない場合
	if err == gorm.ErrRecordNotFound {
		baseChunk := BaseChunk{}

		// チャンクが存在しない場合 (初めて入るチャンクの場合)
		// ベースチャンクからデータを取得
		err = dbconn.Where(BaseChunk{
			RegionID: game.RegionID,
			GridID:   gridId,
		}).First(&baseChunk).Error
		
		if err != nil {
			return GameChunk{}, err
		}

		// チャンクのIDを生成する
		chunkId,_ := utils.Genid()

		// 新く作るチャンクのデータ
		chunkData := GameChunk{
			ChunkID:  "chunkId-" + chunkId,
			GameID:   game.GameID,
			ImageID:  "",
			OwnerID:  "",
			StartLat: baseChunk.StartLat,
			StartLon: baseChunk.StartLon,
			EndLat:   baseChunk.EndLat,
			EndLon:   baseChunk.EndLon,
			GridID:   gridId,
			Level:    0,
		}

		// ベースチャンクを元にゲームチャンクを作成する
		err = dbconn.Create(&chunkData).Error

		// エラー処理
		if err != nil {
			return GameChunk{}, err
		}

		// 新く作ったチャンクのデータを返す
		return chunkData, nil
	}

	// エラー処理
	if err != nil {
		return GameChunk{}, err
	}

	return findChunk, nil
}


func DebugBaseChunk() {
}
