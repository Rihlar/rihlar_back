package location

import (
	"context"
	"errors"
	"gcore/logger"
	"strings"

	"github.com/redis/go-redis/v9"
)

// リージョンが存在するか
func ExistsRegion(regionId string) (bool) {
	// リージョンが存在しない場合
	result := redisConn.Exists(context.Background(), regionId)

	// エラー処理
	if result.Err() != nil {
		return false
	}

	if result.Val() == 1 {
		return true
	}

	return false
}
// 一番近いチャンクを返す
func FindNearChunk(regionId string,Lat float64,Lon float64,Distance float64) (string,error) {
	// リージョンが存在しない場合
	if !ExistsRegion(regionId) {
		return "",nil
	}

	// 一番近いチャンクを返す
	results,err := redisConn.GeoSearchLocation(context.Background(), regionId,&redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  Lon,
			Latitude:   Lat,
			Radius:     Distance,
			RadiusUnit: "m",
			Sort:       "asc",
			CountAny:   true,
		},
		WithCoord:      true,
		WithDist:       true,
		WithHash:       false,
	}).Result()

	// エラー処理
	if err != nil {
		return "",err
	}

	if len(results) > 0 {
		// 全てのチャンクを表示する
		for _, result := range results {
			logger.Println(result)
		}

		// 見つかった場合 チャンクのID (名前を返す)
		return results[0].Name,nil
	}

	return "",errors.New("chunk not found")
}

func SaveChunk(regionID,chunkId string,Lat float64,Lon float64) (error) {
	// grid を保存する
	result := redisConn.GeoAdd(context.Background(), regionID, &redis.GeoLocation{
		Name: chunkId,
		Latitude: Lat,
		Longitude: Lon,
	})

	// エラー処理
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func SaveCircleChunk(regionID,gridId string,Lat float64,Lon float64) (error) {
	// 円の検索用にチャンクを追加する
	result := redisConn.GeoAdd(context.Background(), regionID + "-circle", &redis.GeoLocation{
		Name: gridId,
		Latitude: Lat,
		Longitude: Lon,
	})

	// エラー処理
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// リージョンないのチャンクを保存する
func SaveRegionChunk(regionId string,StartLat float64,StartLon float64,EndLat float64,EndLon float64) (error) {
	// グリッド一覧取得
	grids := GenerateGrid(StartLat,StartLon,EndLat,EndLon,3000)

	for _, grid := range grids {
		// grid を保存する
		err := SaveChunk(regionId,grid.ID,grid.TopLeft.Lat,grid.TopLeft.Lon)

		// エラー処理
		if err != nil {
			return err
		}
	}

	return nil
}

// 接しているチャンクを取得 (円形に検索)
func FindContactChunk(regionId string,Lat float64,Lon float64,Distance float64) ([]string,error) {
	// 円形に近いチャンクを探す
	results,err := redisConn.GeoSearchLocation(context.Background(), regionId + "-circle",&redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  Lon,
			Latitude:   Lat,
			Radius:     Distance,
			RadiusUnit: "m",
			Sort:       "asc",
			CountAny:   true,
		},
		WithCoord:      true,
		WithDist:       true,
		WithHash:       false,
	}).Result()

	// エラー処理
	if err != nil {
		return []string{},err
	}

	returnDatas := []string{}

	if len(results) > 0 {
		// 全てのチャンクを表示する
		for _, result := range results {
			// 後ろの記号を削除する (0_108|left_bottom -> 0_108)
			returnDatas = append(returnDatas,strings.SplitN(result.Name,"|",2)[0])
		}
	}

	return returnDatas,nil
}