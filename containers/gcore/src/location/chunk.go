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
			
			// グリッドの中心点を計算
			center := CalculateCenter(topLeft, bottomRight)
			
			// グリッドIDを生成（行_列の形式）
			gridID := fmt.Sprintf("%d_%d", rowIndex, colIndex)
			
			// Gridを作成してリストに追加
			grid := Grid{
				ID:          gridID,
				TopLeft:     topLeft,
				BottomRight: bottomRight,
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
			
			// 中心座標を計算
			center := CalculateCenter(topLeft, bottomRight)
			
			// グリッドIDを生成
			gridID := fmt.Sprintf("%d_%d", rowIndex, colIndex)
			
			grid := Grid{
				ID:          gridID,
				TopLeft:     topLeft,
				BottomRight: bottomRight,
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

// FindNearestGridByCenter 指定した座標に最も近い中心点を持つグリッドを検索
func FindNearestGridByCenter(grids []Grid, lat, lon float64) (*Grid, float64) {
	if len(grids) == 0 {
		return nil, 0
	}
	
	var nearestGrid *Grid
	minDistance := math.Inf(1)
	
	for i := range grids {
		distance := HaversineDistance(lat, lon, grids[i].Center.Lat, grids[i].Center.Lon)
		if distance < minDistance {
			minDistance = distance
			nearestGrid = &grids[i]
		}
	}
	
	return nearestGrid, minDistance
}

func PrintGrid(grids []Grid) {
	fmt.Printf("グリッド一覧 (総数: %d):\n", len(grids))
	fmt.Println("番号 | 左上座標              | 右下座標              | 中心座標              | 幅(m)   | 高さ(m)")
	fmt.Println(strings.Repeat("-", 100))
	
	for i, grid := range grids {
		// 各グリッドの幅と高さを計算
		width := HaversineDistance(grid.TopLeft.Lat, grid.TopLeft.Lon, grid.TopLeft.Lat, grid.BottomRight.Lon)
		height := HaversineDistance(grid.TopLeft.Lat, grid.TopLeft.Lon, grid.BottomRight.Lat, grid.TopLeft.Lon)
		
		fmt.Printf("%4d | (%.6f,%.6f) | (%.6f,%.6f) | (%.6f,%.6f) | %7.1f | %7.1f\n",
			i, grid.TopLeft.Lat, grid.TopLeft.Lon,
			grid.BottomRight.Lat, grid.BottomRight.Lon,
			grid.Center.Lat, grid.Center.Lon,
			width, height)
		
		// 最初の10個まで表示
		if i >= 9 {
			fmt.Printf("... (残り%d個のグリッド)\n", len(grids)-i-1)
			break
		}
	}
}

func main() {
	// 例1：基本的な使用例
	fmt.Println("=== 基本的な座標移動の例 ===")
	baseLat := 35.681236
	baseLon := 139.767125
	meterEast := 1000.0
	meterSouth := 500.0
	
	fmt.Printf("基準点: 緯度=%.6f, 経度=%.6f\n", baseLat, baseLon)
	fmt.Printf("移動距離: 東へ%.0fm, 南へ%.0fm\n", meterEast, meterSouth)
	
	result := GetBottomRightCornerPrecise(baseLat, baseLon, meterEast, meterSouth)
	fmt.Printf("右下の座標: 緯度=%.6f, 経度=%.6f\n", result.Lat, result.Lon)
	
	// 検証
	distanceEast := HaversineDistance(baseLat, baseLon, baseLat, result.Lon)
	distanceSouth := HaversineDistance(baseLat, baseLon, result.Lat, baseLon)
	fmt.Printf("検証 - 東方向距離: %.2fm, 南方向距離: %.2fm\n", distanceEast, distanceSouth)
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	
	// 例2：グリッド生成の例
	fmt.Println("=== グリッド生成の例 ===")
	
	// 東京駅周辺の小さなエリア（約1km x 1km）
	topLeftLat := 35.685000
	topLeftLon := 139.763000
	bottomRightLat := 35.677000
	bottomRightLon := 139.771000
	gridMeter := 300.0 // 300m間隔
	
	fmt.Printf("左上座標: (%.6f, %.6f)\n", topLeftLat, topLeftLon)
	fmt.Printf("右下座標: (%.6f, %.6f)\n", bottomRightLat, bottomRightLon)
	fmt.Printf("グリッド間隔: %.0fm\n\n", gridMeter)
	
	// 精密計算版でグリッド生成
	grids := GenerateGrid(topLeftLat, topLeftLon, bottomRightLat, bottomRightLon, gridMeter)
	fmt.Printf("精密計算版 - 生成されたグリッド数: %d個\n", len(grids))
	
	// 簡易計算版でグリッド生成
	gridsSimple := GenerateGridSimple(topLeftLat, topLeftLon, bottomRightLat, bottomRightLon, gridMeter)
	fmt.Printf("簡易計算版 - 生成されたグリッド数: %d個\n\n", len(gridsSimple))
	
	// グリッドの詳細を表示
	fmt.Println("精密計算版グリッドの詳細:")
	PrintGrid(grids)
	
	// 実際の距離検証（最初のいくつかのグリッドで）
	if len(grids) > 0 {
		fmt.Println("\n距離検証（最初のグリッド）:")
		firstGrid := grids[0]
		
		// グリッドの幅（水平方向）
		width := HaversineDistance(
			firstGrid.TopLeft.Lat, firstGrid.TopLeft.Lon,
			firstGrid.TopLeft.Lat, firstGrid.BottomRight.Lon)
		
		// グリッドの高さ（垂直方向）
		height := HaversineDistance(
			firstGrid.TopLeft.Lat, firstGrid.TopLeft.Lon,
			firstGrid.BottomRight.Lat, firstGrid.TopLeft.Lon)
		
		fmt.Printf("グリッド0の実際のサイズ: 幅=%.2fm, 高さ=%.2fm\n", width, height)
		fmt.Printf("指定サイズとの差異: 幅=%.2fm, 高さ=%.2fm\n", 
			math.Abs(width-gridMeter), math.Abs(height-gridMeter))
			
		// 中心点から各角への距離を検証
		centerToTopLeft := HaversineDistance(
			firstGrid.Center.Lat, firstGrid.Center.Lon,
			firstGrid.TopLeft.Lat, firstGrid.TopLeft.Lon)
		centerToBottomRight := HaversineDistance(
			firstGrid.Center.Lat, firstGrid.Center.Lon,
			firstGrid.BottomRight.Lat, firstGrid.BottomRight.Lon)
		
		fmt.Printf("中心点から左上までの距離: %.2fm\n", centerToTopLeft)
		fmt.Printf("中心点から右下までの距離: %.2fm\n", centerToBottomRight)
	}
	
	// 各グリッドの座標例
	fmt.Println("\n最初の3つのグリッドの詳細:")
	for i := 0; i < min(3, len(grids)); i++ {
		fmt.Printf("グリッド%s:\n", grids[i].ID)
		fmt.Printf("  左上: (%.6f, %.6f)\n", grids[i].TopLeft.Lat, grids[i].TopLeft.Lon)
		fmt.Printf("  右下: (%.6f, %.6f)\n", grids[i].BottomRight.Lat, grids[i].BottomRight.Lon)
		fmt.Printf("  中心: (%.6f, %.6f)\n", grids[i].Center.Lat, grids[i].Center.Lon)
		
		// 対角線の距離
		diagonal := HaversineDistance(
			grids[i].TopLeft.Lat, grids[i].TopLeft.Lon,
			grids[i].BottomRight.Lat, grids[i].BottomRight.Lon)
		fmt.Printf("  対角線距離: %.2fm\n", diagonal)
		fmt.Println()
	}
	
	fmt.Println("\n" + strings.Repeat("=", 60))
	
	// 例3：座標からグリッド検索の例
	fmt.Println("=== 座標からグリッド検索の例 ===")
	
	// テスト用の座標（東京駅周辺の任意の点）
	testLat := 35.683000
	testLon := 139.765000
	
	fmt.Printf("検索対象座標: (%.6f, %.6f)\n", testLat, testLon)
	
	// グリッドIDを検索
	if gridID, found := FindGridByCoordinate(grids, testLat, testLon); found {
		fmt.Printf("該当グリッドID: %s\n", gridID)
		
		// グリッドの詳細情報を取得
		if grid, exists := GetGridByID(grids, gridID); exists {
			fmt.Printf("グリッドの詳細:\n")
			fmt.Printf("  ID: %s\n", grid.ID)
			fmt.Printf("  左上: (%.6f, %.6f)\n", grid.TopLeft.Lat, grid.TopLeft.Lon)
			fmt.Printf("  右下: (%.6f, %.6f)\n", grid.BottomRight.Lat, grid.BottomRight.Lon)
			fmt.Printf("  中心: (%.6f, %.6f)\n", grid.Center.Lat, grid.Center.Lon)
			
			// 検索点から中心点までの距離
			distanceToCenter := HaversineDistance(testLat, testLon, grid.Center.Lat, grid.Center.Lon)
			fmt.Printf("  検索点から中心点までの距離: %.2fm\n", distanceToCenter)
			
			// 隣接グリッドを取得
			neighbors := GetNeighborGrids(grids, gridID)
			fmt.Printf("  隣接グリッド: %v\n", neighbors)
		}
	} else {
		fmt.Println("該当するグリッドが見つかりませんでした")
	}
	
	// 中心点による最近接検索の例
	fmt.Println("\n=== 中心点による最近接検索 ===")
	if nearestGrid, distance := FindNearestGridByCenter(grids, testLat, testLon); nearestGrid != nil {
		fmt.Printf("最も近い中心点を持つグリッド: %s\n", nearestGrid.ID)
		fmt.Printf("中心点座標: (%.6f, %.6f)\n", nearestGrid.Center.Lat, nearestGrid.Center.Lon)
		fmt.Printf("距離: %.2fm\n", distance)
	}
	
	// 複数の座標でテスト
	testCoords := []LatLon{
		{35.684000, 139.764000}, // 左上寄り
		{35.680000, 139.768000}, // 中央寄り
		{35.678000, 139.770000}, // 右下寄り
		{35.690000, 139.760000}, // 範囲外
	}
	
	fmt.Println("\n複数座標での検索テスト:")
	for i, coord := range testCoords {
		if gridID, found := FindGridByCoordinate(grids, coord.Lat, coord.Lon); found {
			if grid, exists := GetGridByID(grids, gridID); exists {
				fmt.Printf("座標%d (%.6f, %.6f) -> グリッド %s (中心: %.6f, %.6f)\n", 
					i+1, coord.Lat, coord.Lon, gridID, grid.Center.Lat, grid.Center.Lon)
			}
		} else {
			fmt.Printf("座標%d (%.6f, %.6f) -> 範囲外\n", 
				i+1, coord.Lat, coord.Lon)
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