package models

import (
	"gcore/location"
	"gcore/logger"
	"time"
)

func DebugPerformance() {

}

// TODO デバッグ用 ゲーム用のリージョンを作成する関数
func PerformanceFillRegion(region Region) error {
	// 時間を計測する
	start := time.Now()

	// 左上の座標
	topLeft := location.LatLng{
		Lat: region.StartLat,
		Lng: region.StartLon,
	}

	// 右下の座標
	bottomRight := location.LatLng{
		Lat: region.EndLat,
		Lng: region.EndLon,
	}

	grid, err := location.NewRegionGridInfo(topLeft, bottomRight, GridMeter)

	// エラー処理zq
	if err != nil {
		return err
	}

	cell,err := grid.GetGridCell(location.LatLng{
		Lat: 35.7365646305852,
		Lng: 139.80954685258914,
	})

	// エラー処理
	if err != nil {
		return err
	}

	logger.Println("cell", cell)

	// 300m以内のgrid取得
	grids,err := grid.GetGridsInRadius(location.LatLng{
		Lat: 35.7365646305852,
		Lng: 139.80954685258914,
	}, 50)

	// エラー処理
	if err != nil {
		return err
	}

	logger.Println("grids", grids)

	// 経過時間を表示する
	logger.Println("PerformanceFillRegion", time.Since(start))

	return nil
}
