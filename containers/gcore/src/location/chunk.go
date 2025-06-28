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
	BottomRight LatLon // 右下の座標
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
			
			// グリッドIDを生成（行_列の形式）
			gridID := fmt.Sprintf("%d_%d", rowIndex, colIndex)
			
			// Gridを作成してリストに追加
			grid := Grid{
				ID:          gridID,
				TopLeft:     topLeft,
				BottomRight: bottomRight,
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
			
			// グリッドIDを生成
			gridID := fmt.Sprintf("%d_%d", rowIndex, colIndex)
			
			grid := Grid{
				ID:          gridID,
				TopLeft:     topLeft,
				BottomRight: bottomRight,
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

// FindGridByCoordinate 指定した緯度経度が含まれるGridのIDを返す
func FindGridByCoordinate(grids []Grid, lat, lon float64) (string, bool) {
	for _, grid := range grids {
		if IsPointInGrid(grid, lat, lon) {
			return grid.ID, true
		}
	}
	return "", false
}

// FindGridByCoordinateWithDetails 指定した緯度経度が含まれるGridの詳細情報を返す
func FindGridByCoordinateWithDetails(grids []Grid, lat, lon float64) (*Grid, bool) {
	for _, grid := range grids {
		if IsPointInGrid(grid, lat, lon) {
			return &grid, true
		}
	}
	return nil, false
}

// IsPointInGrid 指定した緯度経度がGridの範囲内にあるかチェックする
func IsPointInGrid(grid Grid, lat, lon float64) bool {
	// 緯度チェック（北が大きい値）
	if lat > grid.TopLeft.Lat || lat < grid.BottomRight.Lat {
		return false
	}
	
	// 経度チェック（東が大きい値）
	if lon < grid.TopLeft.Lon || lon > grid.BottomRight.Lon {
		return false
	}
	
	return true
}

// GetGridByID 指定したIDのGridを取得する
func GetGridByID(grids []Grid, id string) (*Grid, bool) {
	for _, grid := range grids {
		if grid.ID == id {
			return &grid, true
		}
	}
	return nil, false
}

// GetNeighborGrids 指定したGridの隣接するGridのIDリストを返す
func GetNeighborGrids(grids []Grid, targetID string) []string {
	var neighbors []string
	
	// IDから行列番号を取得
	var targetRow, targetCol int
	fmt.Sscanf(targetID, "%d_%d", &targetRow, &targetCol)
	
	// 8方向の隣接グリッドをチェック
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // 上行
		{0, -1}, {0, 1},           // 同行
		{1, -1}, {1, 0}, {1, 1},   // 下行
	}
	
	for _, dir := range directions {
		neighborRow := targetRow + dir[0]
		neighborCol := targetCol + dir[1]
		neighborID := fmt.Sprintf("%d_%d", neighborRow, neighborCol)
		
		// 隣接グリッドが実際に存在するかチェック
		if _, exists := GetGridByID(grids, neighborID); exists {
			neighbors = append(neighbors, neighborID)
		}
	}
	
	return neighbors
}
func PrintGrid(grids []Grid) {
	fmt.Printf("グリッド一覧 (総数: %d):\n", len(grids))
	fmt.Println("番号 | 左上座標              | 右下座標              | 幅(m)   | 高さ(m)")
	fmt.Println(strings.Repeat("-", 80))
	
	for i, grid := range grids {
		// 各グリッドの幅と高さを計算
		width := HaversineDistance(grid.TopLeft.Lat, grid.TopLeft.Lon, grid.TopLeft.Lat, grid.BottomRight.Lon)
		height := HaversineDistance(grid.TopLeft.Lat, grid.TopLeft.Lon, grid.BottomRight.Lat, grid.TopLeft.Lon)
		
		fmt.Printf("%4d | (%.6f,%.6f) | (%.6f,%.6f) | %7.1f | %7.1f\n",
			i, grid.TopLeft.Lat, grid.TopLeft.Lon,
			grid.BottomRight.Lat, grid.BottomRight.Lon,
			width, height)
		
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