package location

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CircleData struct {
	Center   LatLng  //中心
	Radius   float64 //半径
	CircleID string  //円のID
	UserID string	 //メンバーID (ユーザーID)
}

func CacheCircle(data CircleData) error {
	// データをキャッシュする
	result := redisConn.GeoAdd(context.Background(),data.UserID, &redis.GeoLocation{
		Name:      data.CircleID,
		Longitude: data.Center.Lng,
		Latitude:  data.Center.Lat,
	})

	// エラー処理
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// キャッシュから削除する関数
func DeleteCircle(data CircleData) error {
	// データをキャッシュする
	result := redisConn.ZRem(context.Background(),data.UserID, data.CircleID)

	// エラー処理
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

type SearchResult struct {
	CircleID string  //円のID
	Distance float64 //距離
}

// 近い円を取得する
func GetNearCircle(UserID string,center LatLng, radius float64) ([]SearchResult,error) {
	// データを取得する
	results := redisConn.GeoSearchLocation(context.Background(), UserID, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  center.Lng,
			Latitude:   center.Lat,
			Radius:     radius,
			RadiusUnit: "m",
			Sort:       "ASC",
		},
		WithCoord:      false,
		WithDist:       true,
		WithHash:       false,
	})

	// エラー処理
	if results.Err() != nil {
		return nil, results.Err()
	}

	// 返すデータ
	returnDatas := []SearchResult{}

	for _, val := range results.Val() {
		// 結果に追加
		returnDatas = append(returnDatas, SearchResult{
			CircleID: val.Name,
			Distance: val.Dist,
		})
	}
	
	return returnDatas, nil
}