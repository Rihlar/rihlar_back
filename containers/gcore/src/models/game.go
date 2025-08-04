package models

import (
	"gcore/location"
	"gcore/logger"
	"gcore/utils"
	"time"

	"gorm.io/gorm"
)

// テーブル定義
type Game struct {
	GameID    string    `gorm:"primaryKey;size:50" json:"gameID"`                                                               // ゲームID
	StartTime time.Time `gorm:"DATETIME;not null" json:"startTime"`                                                             // ゲーム開始時間
	EndTime   time.Time `gorm:"DATETIME;not null" json:"endTime"`                                                               // ゲーム終了時間
	Flag      int       `gorm:"not null" json:"flag"`                                                                           // ゲームユニット 0:個人戦、1:チーム戦
	Type      int       `gorm:"not null" json:"type"`                                                                           // ゲームタイプ	0:system、1:admin
	Teams     []Team    `gorm:"foreignKey:GameID;references:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams"` // 参加チーム
	Status    int       `gorm:"not null" json:"status"`                                                                         // ゲームステータス　0:開始前、1:開催中、2:終了済
	RegionID  string    `gorm:"varchar(50)" json:"regionID"`                                                                    // ゲーム開催地域
}

// テーブル名
func (Game) TableName() string {
	return "games"
}

// TODO デバッグ用 チームを追加する
func (game *Game) AddTeam(team Team) error {
	return dbconn.Model(game).Association("Teams").Append(&team)
}

// ゲームを取得するエンドポイント
func GetGame(gameid string) (Game, error) {
	var game Game

	// 取得する
	result := dbconn.Where(&Game{
		GameID: gameid,
	}).First(&game)

	return game, result.Error
}

// ゲームを保存するエンドポイント
func SaveGame(game Game) error {
	return dbconn.Save(&game).Error
}

// 緯度経度からチャンクを取得する
func (game *Game) GetChunkByLatLon(lat, lon float64) (GameChunk, error) {
	// リージョンの情報を取得
	region, err := GetRegionByID(game.RegionID)

	// エラー処理
	if err != nil {
		return GameChunk{}, err
	}

	// 自分が入っているチャンクを取得
	// リージョンから情報を生成
	gridData,err := location.NewRegionGridInfo(location.LatLng{
		Lat: region.StartLat,
		Lng: region.StartLon,
	},location.LatLng{
		Lat: region.EndLat,
		Lng: region.EndLon,
	},GridMeter)

	// エラー処理
	if err != nil {
		return GameChunk{}, err
	}

	// Grid から近くのチャンクを取得
	inGridData, err := gridData.GetGridCell(location.LatLng{
		Lat: lat,
		Lng: lon,
	})

	// エラー処理
	if err != nil {
		return GameChunk{}, err
	}

	findChunk := GameChunk{}
	// ゲームからチャンクを取得する
	err = dbconn.Where(GameChunk{
		GameID:   game.GameID,
		StartLat: inGridData.Bounds.TopLeft.Lat,
		StartLon: inGridData.Bounds.TopLeft.Lng,
		EndLat:   inGridData.Bounds.BottomRight.Lat,
		EndLon:   inGridData.Bounds.BottomRight.Lng,
	}).First(&findChunk).Error

	// 存在しない場合
	if err == gorm.ErrRecordNotFound {
		// チャンクのIDを生成する
		chunkId,_ := utils.Genid()

		// 新く作るチャンクのデータ
		chunkData := GameChunk{
			ChunkID:  "chunkId-" + chunkId,
			GameID:   game.GameID,
			ImageID:  "",
			OwnerID:  "",
			StartLat: inGridData.Bounds.TopLeft.Lat,
			StartLon: inGridData.Bounds.TopLeft.Lng,
			EndLat:   inGridData.Bounds.BottomRight.Lat,
			EndLon:   inGridData.Bounds.BottomRight.Lng,
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

// 円形にチャンクを取得する
func (game *Game) GetCircleChunkByLatLon(lat, lon, radius float64) ([]GameChunk, error) {
	// リージョンを取得
	region, err := GetRegionByID(game.RegionID)

	// エラー処理	
	if err != nil {
		return []GameChunk{}, err
	}

	// 自分が入っているチャンクを取得
	// リージョンから情報を生成
	gridData,err := location.NewRegionGridInfo(location.LatLng{
		Lat: region.StartLat,
		Lng: region.StartLon,
	},location.LatLng{
		Lat: region.EndLat,
		Lng: region.EndLon,
	},GridMeter)

	// エラー処理
	if err != nil {
		return []GameChunk{}, err
	}

	// Grid から近くのチャンクを取得
	inGridDatas, err := gridData.GetGridsInRadius(location.LatLng{
		Lat: lat,
		Lng: lon,
	}, radius)

	// エラー処理
	if err != nil {
		return []GameChunk{}, err
	}

	// 取得したチャンクを格納する
	returnChunks := []GameChunk{}

	for _, inGrid := range inGridDatas {
		// ゲームチャンクを取得する
		findChunk := GameChunk{}
		// ゲームからチャンクを取得する
		err = dbconn.Where(GameChunk{
			GameID:   game.GameID,
			StartLat: inGrid.Bounds.TopLeft.Lat,
			StartLon: inGrid.Bounds.TopLeft.Lng,
			EndLat:   inGrid.Bounds.BottomRight.Lat,
			EndLon:   inGrid.Bounds.BottomRight.Lng,
		}).First(&findChunk).Error

		// 存在しない場合
		if err == gorm.ErrRecordNotFound {
			// チャンクのIDを生成する
			chunkId,_ := utils.Genid()

			// 新く作るチャンクのデータ
			chunkData := GameChunk{
				ChunkID:  "chunkId-" + chunkId,
				GameID:   game.GameID,
				ImageID:  "",
				OwnerID:  "",
				StartLat: inGrid.Bounds.TopLeft.Lat,
				StartLon: inGrid.Bounds.TopLeft.Lng,
				EndLat:   inGrid.Bounds.BottomRight.Lat,
				EndLon:   inGrid.Bounds.BottomRight.Lng,
				Level:    0,
			}

			// ベースチャンクを元にゲームチャンクを作成する
			err = dbconn.Create(&chunkData).Error

			// エラー処理
			if err != nil {
				return []GameChunk{}, err
			}

			// 新く作ったチャンクのデータを返す
			returnChunks = append(returnChunks, chunkData)
			continue
		}

		// エラー処理
		if err != nil {
			return []GameChunk{}, err
		}

		returnChunks = append(returnChunks, findChunk)
	}

	return returnChunks, nil
}

// ランキング上位取得
func (game *Game) GetRanking(maxRank int) ([]Team, error) {
	var rankings []Team

	result := dbconn.Debug().
		Where(Team{
			GameID:    game.GameID,
		}).
		Order("points DESC").
		Limit(maxRank).
		Find(&rankings)

	if result.Error != nil {
		logger.PrintErr("ランキング上位取得エラー", result.Error)
		return nil, result.Error
	}

	return rankings, nil
}

// デバック用
func DebugGame() {
	// テスト用のゲームデータを入れる
	CreateTestGames()
}

func CreateTestGames() {
	// システムゲームを作成する
	for _, sysgame := range SysGameIDs {
		CreateGame(Game{
			GameID:    sysgame,
			StartTime: time.Now(),
			EndTime:   time.Now().AddDate(0,0,20),
			Flag:      0,
			Type:      0,
			Status:    1,
			RegionID:  RegionId,
		})
	}

	// admin ゲームを作成する
	CreateGame(Game{
		GameID:    AdminGameId1,
		StartTime: time.Now(),
		EndTime:   time.Now().AddDate(0,0,20),
		Flag:      0,
		Type:      1,
		Status:    1,
		RegionID:  RegionId,
	})

	// ユーザーをシステムゲームに追加していく
	for index, sysgame := range SysGameIDs {
		// チームIDを生成する
		teamId,_ := utils.Genid()

		DebugAddMember(sysgame, "teamid-" + teamId, UserIDs[index])

		// チームIDを生成する
		teamId2,_ := utils.Genid()

		// admin ゲームに追加していく
		DebugAddMember(AdminGameId1, "teamid-" + teamId2, UserIDs[index])
	}
}

func CreateGame(game Game) error {
	return dbconn.Create(&game).Error
}

// 一人のゲーム追加をデバッグする
func DebugAddMember(gameID string, teamID string, userID string) error {
	// ゲームを取得する
	game, err := GetGame(gameID)

	// エラー処理
	if err != nil {
		logger.PrintErr("ゲーム取得エラー", err)
		return err
	}

	// ゲームにチームを追加する
	err = game.AddTeam(Team{
		TeamID: teamID,
		Points: 0,
	})

	// エラー処理
	if err != nil {
		return err
	}

	// チームを取得
	team, err := game.GetTeam(teamID)

	// エラー処理
	if err != nil {
		return err
	}

	// チームにメンバーを追加する
	err = team.AddMember(Member{
		UserID: userID,
		Points: 0,
	})

	// エラー処理
	if err != nil {
		logger.PrintErr("メンバー追加エラー", err)
		return err
	}

	return nil
}

