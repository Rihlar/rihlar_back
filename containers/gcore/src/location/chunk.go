package location

import (
	"fmt"
	"math"
	"strings"
)

// 地球の半径（メートル）
const EarthRadius = 6378137.0

// LatLon 緯度経度を表す構造体
type LatLon struct {
	Lat float64 // 緯度
	Lon float64 // 経度
}

// Grid グリッドの1つのセルを表す構造体
type Grid struct {
	ID          string // グリッドID
	TopLeft     LatLon // 左上の座標
	TopRight    LatLon // 右上の座標
	BottomLeft  LatLon // 左下の座標
	BottomRight LatLon // 右下の座標
	Center      LatLon // 中心座標
}

// CalculateCenter グリッドの中心点を計算する
func CalculateCenter(topLeft, bottomRight LatLon) LatLon {
	// 単純に緯度経度の中点を計算
	centerLat := (topLeft.Lat + bottomRight.Lat) / 2.0
	centerLon := (topLeft.Lon + bottomRight.Lon) / 2.0
	
	return LatLon{
		Lat: centerLat,
		Lon: centerLon,
	}
}

// CalculateCenterPrecise より精密な中心点計算（測地線を考慮）
func CalculateCenterPrecise(topLeft, bottomRight LatLon) LatLon {
	// 緯度経度をラジアンに変換
	lat1 := topLeft.Lat * math.Pi / 180.0
	lon1 := topLeft.Lon * math.Pi / 180.0
	lat2 := bottomRight.Lat * math.Pi / 180.0
	lon2 := bottomRight.Lon * math.Pi / 180.0
	
	// 中点の計算（球面上の測地線中点）
	deltaLon := lon2 - lon1
	
	bx := math.Cos(lat2) * math.Cos(deltaLon)
	by := math.Cos(lat2) * math.Sin(deltaLon)
	
	lat3 := math.Atan2(math.Sin(lat1)+math.Sin(lat2),
		math.Sqrt((math.Cos(lat1)+bx)*(math.Cos(lat1)+bx)+by*by))
	lon3 := lon1 + math.Atan2(by, math.Cos(lat1)+bx)
	
	// ラジアンから度に変換
	centerLat := lat3 * 180.0 / math.Pi
	centerLon := lon3 * 180.0 / math.Pi
	
	return LatLon{
		Lat: centerLat,
		Lon: centerLon,
	}
}

// Calculate4Corners 左上と右下の座標から4つの角の座標を計算する
func Calculate4Corners(topLeft, bottomRight LatLon) (LatLon, LatLon, LatLon, LatLon) {
	// 右上座標：左上の緯度、右下の経度
	topRight := LatLon{
		Lat: topLeft.Lat,
		Lon: bottomRight.Lon,
	}
	
	// 左下座標：右下の緯度、左上の経度
	bottomLeft := LatLon{
		Lat: bottomRight.Lat,
		Lon: topLeft.Lon,
	}
	
	return topLeft, topRight, bottomLeft, bottomRight
}

// GetBottomRightCorner 指定した緯度経度とメートル数から、グリッドの右下の緯度経度を計算する
func GetBottomRightCorner(baseLat, baseLon float64, meterEast, meterSouth float64) LatLon {
	// 緯度1度あたりのメートル数（約111,111m）
	latPerMeter := 1.0 / 111111.0
	
	// 経度1度あたりのメートル数（緯度により変化）
	// cos(緯度) * 111,111m
	lonPerMeter := 1.0 / (111111.0 * math.Cos(baseLat*math.Pi/180.0))
	
	// 右（東）へ移動した経度
	newLon := baseLon + (meterEast * lonPerMeter)
	
	// 下（南）へ移動した緯度
	newLat := baseLat - (meterSouth * latPerMeter)
	
	return LatLon{
		Lat: newLat,
		Lon: newLon,
	}
}

// より精密な計算を行う関数（Haversine公式を使用）
func GetBottomRightCornerPrecise(baseLat, baseLon float64, meterEast, meterSouth float64) LatLon {
	// 基準点をラジアンに変換
	baseLatRad := baseLat * math.Pi / 180.0
	baseLonRad := baseLon * math.Pi / 180.0
	
	// 東方向への移動
	deltaLonRad := meterEast / (EarthRadius * math.Cos(baseLatRad))
	newLonRad := baseLonRad + deltaLonRad
	
	// 南方向への移動
	deltaLatRad := -meterSouth / EarthRadius // 南は負の方向
	newLatRad := baseLatRad + deltaLatRad
	
	// ラジアンから度に変換
	newLat := newLatRad * 180.0 / math.Pi
	newLon := newLonRad * 180.0 / math.Pi
	
	return LatLon{
		Lat: newLat,
		Lon: newLon,
	}
}

// 距離を計算する関数（検証用）
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// 度をラジアンに変換
	lat1Rad := lat1 * math.Pi / 180.0
	lon1Rad := lon1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	lon2Rad := lon2 * math.Pi / 180.0
	
	// 差分を計算
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad
	
	// Haversine公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return EarthRadius * c
}

// GenerateGrid 左上と右下の座標、グリッドのメートル間隔からグリッドのリストを生成する
func GenerateGrid(topLeftLat, topLeftLon, bottomRightLat, bottomRightLon float64, gridMeter float64) []Grid {
	var grids []Grid
	
	// 北から南への緯度方向の処理
	currentTopLat := topLeftLat
	rowIndex := 0
	
	for currentTopLat > bottomRightLat {
		currentLeftLon := topLeftLon
		colIndex := 0
		
		// 西から東への経度方向の処理
		for currentLeftLon < bottomRightLon {
			// 現在のセルの左上座標
			topLeft := LatLon{
				Lat: currentTopLat,
				Lon: currentLeftLon,
			}
			
			// 現在のセルの右下座標を計算
			bottomRight := GetBottomRightCornerPrecise(currentTopLat, currentLeftLon, gridMeter, gridMeter)
			
			// 境界チェック：右下が指定範囲を超えないようにクリップ
			if bottomRight.Lat < bottomRightLat {
				bottomRight.Lat = bottomRightLat
			}
			if bottomRight.Lon > bottomRightLon {
				bottomRight.Lon = bottomRightLon
			}
			
			// 4つの角の座標を計算
			tl, tr, bl, br := Calculate4Corners(topLeft, bottomRight)
			
			// グリッドの中心点を計算
			center := CalculateCenter(topLeft, bottomRight)
			
			// グリッドIDを生成（行_列の形式）
			gridID := fmt.Sprintf("%d_%d", rowIndex, colIndex)
			
			// Gridを作成してリストに追加
			grid := Grid{
				ID:          gridID,
				TopLeft:     tl,
				TopRight:    tr,
				BottomLeft:  bl,
				BottomRight: br,
				Center:      center,
			}
			grids = append(grids, grid)
			
			// 東へgridMeter分移動した次の経度を計算
			nextPoint := GetBottomRightCornerPrecise(currentTopLat, currentLeftLon, gridMeter, 0)
			currentLeftLon = nextPoint.Lon
			colIndex++
			
			// 境界を超えた場合は終了
			if currentLeftLon >= bottomRightLon {
				break
			}
		}
		
		// 南へgridMeter分移動した次の緯度を計算
		nextPoint := GetBottomRightCornerPrecise(currentTopLat, topLeftLon, 0, gridMeter)
		currentTopLat = nextPoint.Lat
		rowIndex++
		
		// 境界を超えた場合は終了
		if currentTopLat <= bottomRightLat {
			break
		}
	}
	
	return grids
}

// GenerateGridSimple 簡易版のグリッド生成関数
func GenerateGridSimple(topLeftLat, topLeftLon, bottomRightLat, bottomRightLon float64, gridMeter float64) []Grid {
	var grids []Grid
	
	// 緯度1度あたりのメートル数
	latPerMeter := 1.0 / 111111.0
	// 経度1度あたりのメートル数（中央緯度で計算）
	centerLat := (topLeftLat + bottomRightLat) / 2
	lonPerMeter := 1.0 / (111111.0 * math.Cos(centerLat*math.Pi/180.0))
	
	// グリッド間隔を度に変換
	latStep := gridMeter * latPerMeter
	lonStep := gridMeter * lonPerMeter
	
	currentTopLat := topLeftLat
	rowIndex := 0
	
	for currentTopLat > bottomRightLat {
		currentLeftLon := topLeftLon
		colIndex := 0
		
		for currentLeftLon < bottomRightLon {
			// 左上座標
			topLeft := LatLon{
				Lat: currentTopLat,
				Lon: currentLeftLon,
			}
			
			// 右下座標
			bottomRight := LatLon{
				Lat: math.Max(currentTopLat-latStep, bottomRightLat),
				Lon: math.Min(currentLeftLon+lonStep, bottomRightLon),
			}
			
			// 4つの角の座標を計算
			tl, tr, bl, br := Calculate4Corners(topLeft, bottomRight)
			
			// 中心座標を計算
			center := CalculateCenter(topLeft, bottomRight)
			
			// グリッドIDを生成
			gridID := fmt.Sprintf("%d_%d", rowIndex, colIndex)
			
			grid := Grid{
				ID:          gridID,
				TopLeft:     tl,
				TopRight:    tr,
				BottomLeft:  bl,
				BottomRight: br,
				Center:      center,
			}
			grids = append(grids, grid)
			
			currentLeftLon += lonStep
			colIndex++
		}
		
		currentTopLat -= latStep
		rowIndex++
	}
	
	return grids
}

func PrintGrid(grids []Grid) {
	fmt.Printf("グリッド一覧 (総数: %d):\n", len(grids))
	fmt.Println("ID   | 左上座標              | 右上座標              | 左下座標              | 右下座標              | 中心座標")
	fmt.Println(strings.Repeat("-", 130))
	
	for i, grid := range grids {
		fmt.Printf("%s | (%.6f,%.6f) | (%.6f,%.6f) | (%.6f,%.6f) | (%.6f,%.6f) | (%.6f,%.6f)\n",
			grid.ID,
			grid.TopLeft.Lat, grid.TopLeft.Lon,
			grid.TopRight.Lat, grid.TopRight.Lon,
			grid.BottomLeft.Lat, grid.BottomLeft.Lon,
			grid.BottomRight.Lat, grid.BottomRight.Lon,
			grid.Center.Lat, grid.Center.Lon)
		
		// 最初の10個まで表示
		if i >= 9 {
			fmt.Printf("... (残り%d個のグリッド)\n", len(grids)-i-1)
			break
		}
	}
}

// min 関数（Go 1.21未満の場合）
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}