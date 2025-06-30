package location

import (
	"errors"
	"fmt"
	"math"
)

// LatLngは地理座標を表します
type LatLng struct {
	Lat, Lng float64
}

// Pointはメートル単位の投影座標を表します
type Point struct {
	X, Y float64
}

// RegionGridInfoは特定地域のグリッド情報を格納します
type RegionGridInfo struct {
	// グリッドの境界（緯度経度）
	TopLeft     LatLng
	BottomRight LatLng

	// グリッドの設定
	GridSizeMeters float64
	Cols           int
	Rows           int

	// 投影された境界（メートル単位、ウェブメルカトル）
	projTopLeft     Point
	projBottomRight Point

	// パフォーマンス向上のための事前計算済み定数
	earthRadius float64
	degToRad    float64
	radToDeg    float64
	invGridSize float64 // 高速な除算のための 1 / GridSizeMeters
}

// GridCellは座標を持つグリッドセルを表します
type GridCell struct {
	Row, Col int        // グリッド座標（0から始まる）
	Bounds   GridBounds // 地理的な境界
}

// GridBoundsはグリッドセルの地理的な境界を表します
type GridBounds struct {
	TopLeft     LatLng
	BottomRight LatLng
}

// NewRegionGridInfoは新しい地域固有のグリッドシステムを作成します
func NewRegionGridInfo(topLeft, bottomRight LatLng, gridSizeMeters float64) (*RegionGridInfo, error) {
	if gridSizeMeters <= 0 {
		return nil, errors.New("グリッドサイズは正の値でなければなりません")
	}

	// 座標を検証
	if topLeft.Lat <= bottomRight.Lat {
		return nil, errors.New("左上の緯度は右下の緯度より大きくなければなりません")
	}
	if topLeft.Lng >= bottomRight.Lng {
		return nil, errors.New("左上の経度は右下の経度より小さくなければなりません")
	}
	if topLeft.Lat < -85.0511 || topLeft.Lat > 85.0511 || bottomRight.Lat < -85.0511 || bottomRight.Lat > 85.0511 {
		return nil, errors.New("緯度はウェブメルカトルの範囲内（-85.0511から85.0511）でなければなりません")
	}

	earthRadius := 6378137.0 // WGS84地球半径（メートル）
	degToRad := math.Pi / 180.0

	grid := &RegionGridInfo{
		TopLeft:        topLeft,
		BottomRight:    bottomRight,
		GridSizeMeters: gridSizeMeters,
		earthRadius:    earthRadius,
		degToRad:       degToRad,
		radToDeg:       180.0 / math.Pi,
		invGridSize:    1.0 / gridSizeMeters,
	}

	// 境界をウェブメルカトル図法に変換
	grid.projTopLeft = grid.latLngToWebMercator(topLeft)
	grid.projBottomRight = grid.latLngToWebMercator(bottomRight)

	// グリッドの次元を計算
	width := grid.projBottomRight.X - grid.projTopLeft.X
	height := grid.projTopLeft.Y - grid.projBottomRight.Y

	grid.Cols = int(math.Ceil(width / gridSizeMeters))
	grid.Rows = int(math.Ceil(height / gridSizeMeters))

	return grid, nil
}

// latLngToWebMercatorは緯度経度をウェブメルカトル図法（メートル）に変換します
func (g *RegionGridInfo) latLngToWebMercator(coord LatLng) Point {
	// 緯度を有効なウェブメルカトル範囲にクランプ
	lat := math.Max(-85.0511, math.Min(85.0511, coord.Lat))

	// ウェブメルカトルに変換
	x := coord.Lng * g.degToRad * g.earthRadius
	y := math.Log(math.Tan((90+lat)*g.degToRad/2)) * g.earthRadius

	return Point{X: x, Y: y}
}

// webMercatorToLatLngはウェブメルカトル図法を緯度経度に戻します
func (g *RegionGridInfo) webMercatorToLatLng(point Point) LatLng {
	lng := point.X / g.earthRadius * g.radToDeg
	lat := (2*math.Atan(math.Exp(point.Y/g.earthRadius)) - math.Pi/2) * g.radToDeg

	return LatLng{Lat: lat, Lng: lng}
}

// GetGridCellは指定された座標を含むグリッドセルを返します
func (g *RegionGridInfo) GetGridCell(coord LatLng) (*GridCell, error) {
	// 境界チェック
	if coord.Lat > g.TopLeft.Lat || coord.Lat < g.BottomRight.Lat ||
		coord.Lng < g.TopLeft.Lng || coord.Lng > g.BottomRight.Lng {
		return nil, fmt.Errorf("座標 (%.6f, %.6f) はグリッド範囲外です", coord.Lat, coord.Lng)
	}

	// ウェブメルカトルに変換
	mercator := g.latLngToWebMercator(coord)

	// グリッド座標を計算
	col := int((mercator.X - g.projTopLeft.X) * g.invGridSize)
	row := int((g.projTopLeft.Y - mercator.Y) * g.invGridSize)

	// 有効範囲にクランプ（浮動小数点精度のエッジケースを処理）
	if col >= g.Cols {
		col = g.Cols - 1
	}
	if row >= g.Rows {
		row = g.Rows - 1
	}
	if col < 0 {
		col = 0
	}
	if row < 0 {
		row = 0
	}

	// グリッドの境界を計算
	bounds := g.calculateGridBounds(row, col)

	return &GridCell{
		Row:    row,
		Col:    col,
		Bounds: bounds,
	}, nil
}

// GetGridCellFastは境界検証なしでグリッドセルを返します（高速）
func (g *RegionGridInfo) GetGridCellFast(coord LatLng) *GridCell {
	// ウェブメルカトルに変換（安定性のために緯度をクランプ）
	lat := coord.Lat
	if lat > 85.0511 {
		lat = 85.0511
	} else if lat < -85.0511 {
		lat = -85.0511
	}

	x := coord.Lng * g.degToRad * g.earthRadius
	y := math.Log(math.Tan((90+lat)*g.degToRad/2)) * g.earthRadius

	// グリッド座標を計算
	col := int((x - g.projTopLeft.X) * g.invGridSize)
	row := int((g.projTopLeft.Y - y) * g.invGridSize)

	// 有効範囲にクランプ
	if col >= g.Cols {
		col = g.Cols - 1
	} else if col < 0 {
		col = 0
	}
	if row >= g.Rows {
		row = g.Rows - 1
	} else if row < 0 {
		row = 0
	}

	// グリッドの境界を計算
	bounds := g.calculateGridBounds(row, col)

	return &GridCell{
		Row:    row,
		Col:    col,
		Bounds: bounds,
	}
}

// calculateGridBoundsはグリッドセルの地理的境界を計算します
func (g *RegionGridInfo) calculateGridBounds(row, col int) GridBounds {
	// ウェブメルカトル境界を計算
	minX := g.projTopLeft.X + float64(col)*g.GridSizeMeters
	maxX := minX + g.GridSizeMeters
	maxY := g.projTopLeft.Y - float64(row)*g.GridSizeMeters
	minY := maxY - g.GridSizeMeters

	// 緯度経度に戻す
	topLeft := g.webMercatorToLatLng(Point{X: minX, Y: maxY})
	bottomRight := g.webMercatorToLatLng(Point{X: maxX, Y: minY})

	return GridBounds{
		TopLeft:     topLeft,
		BottomRight: bottomRight,
	}
}

// GetGridAtCoordinatesは特定のグリッド座標にあるグリッドセルを返します
func (g *RegionGridInfo) GetGridAtCoordinates(row, col int) (*GridCell, error) {
	if row < 0 || row >= g.Rows || col < 0 || col >= g.Cols {
		return nil, fmt.Errorf("グリッド座標 [%d,%d] は範囲外です (0-%d, 0-%d)", row, col, g.Rows-1, g.Cols-1)
	}

	bounds := g.calculateGridBounds(row, col)

	return &GridCell{
		Row:    row,
		Col:    col,
		Bounds: bounds,
	}, nil
}

// GetNeighboringGridsは隣接するグリッドセル（最大8つ）を返します
func (g *RegionGridInfo) GetNeighboringGrids(row, col int) []*GridCell {
	var neighbors []*GridCell

	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue // 中央のセルはスキップ
			}

			newRow, newCol := row+dr, col+dc
			if newRow >= 0 && newRow < g.Rows && newCol >= 0 && newCol < g.Cols {
				cell, _ := g.GetGridAtCoordinates(newRow, newCol)
				neighbors = append(neighbors, cell)
			}
		}
	}

	return neighbors
}

// GetAllGridsは地域内のすべてのグリッドセルを返します
func (g *RegionGridInfo) GetAllGrids() []*GridCell {
	cells := make([]*GridCell, 0, g.Rows*g.Cols)

	for row := 0; row < g.Rows; row++ {
		for col := 0; col < g.Cols; col++ {
			cell, _ := g.GetGridAtCoordinates(row, col)
			cells = append(cells, cell)
		}
	}

	return cells
}

// GetGridsInRadiusは指定された半径の円と交差するすべてのグリッドセルを返します
func (g *RegionGridInfo) GetGridsInRadius(center LatLng, radiusMeters float64) ([]*GridCell, error) {
	if radiusMeters <= 0 {
		return nil, errors.New("半径は正の値でなければなりません")
	}

	// 中心が領域の近くにあるかを確認します（エッジケースのために少し拡大します）
	expandedRadius := radiusMeters + g.GridSizeMeters
	if !g.isCoordinateNearRegion(center, expandedRadius) {
		return []*GridCell{}, nil // 円が領域と交差しない場合は空のスライスを返します
	}

	// 中心をウェブメルカトル図法に変換します
	centerMercator := g.latLngToWebMercator(center)

	// ウェブメルカトル座標で円のバウンディングボックスを計算します
	minX := centerMercator.X - radiusMeters
	maxX := centerMercator.X + radiusMeters
	minY := centerMercator.Y - radiusMeters
	maxY := centerMercator.Y + radiusMeters

	// バウンディングボックスをグリッド座標に変換します
	minCol := int(math.Floor((minX - g.projTopLeft.X) * g.invGridSize))
	maxCol := int(math.Floor((maxX - g.projTopLeft.X) * g.invGridSize))
	minRow := int(math.Floor((g.projTopLeft.Y - maxY) * g.invGridSize))
	maxRow := int(math.Floor((g.projTopLeft.Y - minY) * g.invGridSize))

	// 有効なグリッド範囲にクランプします
	minCol = int(math.Max(0, float64(minCol)))
	maxCol = int(math.Min(float64(g.Cols-1), float64(maxCol)))
	minRow = int(math.Max(0, float64(minRow)))
	maxRow = int(math.Min(float64(g.Rows-1), float64(maxRow)))

	var result []*GridCell
	radiusSquared := radiusMeters * radiusMeters

	// バウンディングボックス内の各グリッドセルをチェックします
	for row := minRow; row <= maxRow; row++ {
		for col := minCol; col <= maxCol; col++ {
			// ウェブメルカトルでグリッドセルの境界を計算します
			cellMinX := g.projTopLeft.X + float64(col)*g.GridSizeMeters
			cellMaxX := cellMinX + g.GridSizeMeters
			cellMaxY := g.projTopLeft.Y - float64(row)*g.GridSizeMeters
			cellMinY := cellMaxY - g.GridSizeMeters

			// 円がこのグリッドセルの長方形と交差するかどうかを確認します
			if g.circleIntersectsRectangle(centerMercator, radiusSquared, cellMinX, cellMinY, cellMaxX, cellMaxY) {
				bounds := g.calculateGridBounds(row, col)
				result = append(result, &GridCell{
					Row:    row,
					Col:    col,
					Bounds: bounds,
				})
			}
		}
	}

	return result, nil
}

// GetGridsInRadiusFastは境界検証なしで半径内のグリッドセルを返します（高速）
func (g *RegionGridInfo) GetGridsInRadiusFast(center LatLng, radiusMeters float64) []*GridCell {
	// 中心をウェブメルカトルに変換します（緯度クランプ付き）
	lat := center.Lat
	if lat > 85.0511 {
		lat = 85.0511
	} else if lat < -85.0511 {
		lat = -85.0511
	}

	centerX := center.Lng * g.degToRad * g.earthRadius
	centerY := math.Log(math.Tan((90+lat)*g.degToRad/2)) * g.earthRadius

	// バウンディングボックスを計算します
	minX := centerX - radiusMeters
	maxX := centerX + radiusMeters
	minY := centerY - radiusMeters
	maxY := centerY + radiusMeters

	// グリッド座標に変換します
	minCol := int(math.Max(0, math.Floor((minX-g.projTopLeft.X)*g.invGridSize)))
	maxCol := int(math.Min(float64(g.Cols-1), math.Floor((maxX-g.projTopLeft.X)*g.invGridSize)))
	minRow := int(math.Max(0, math.Floor((g.projTopLeft.Y-maxY)*g.invGridSize)))
	maxRow := int(math.Min(float64(g.Rows-1), math.Floor((g.projTopLeft.Y-minY)*g.invGridSize)))

	var result []*GridCell
	radiusSquared := radiusMeters * radiusMeters
	centerMercator := Point{X: centerX, Y: centerY}

	// 各グリッドセルをチェックします
	for row := minRow; row <= maxRow; row++ {
		for col := minCol; col <= maxCol; col++ {
			cellMinX := g.projTopLeft.X + float64(col)*g.GridSizeMeters
			cellMaxX := cellMinX + g.GridSizeMeters
			cellMaxY := g.projTopLeft.Y - float64(row)*g.GridSizeMeters
			cellMinY := cellMaxY - g.GridSizeMeters

			if g.circleIntersectsRectangle(centerMercator, radiusSquared, cellMinX, cellMinY, cellMaxX, cellMaxY) {
				bounds := g.calculateGridBounds(row, col)
				result = append(result, &GridCell{
					Row:    row,
					Col:    col,
					Bounds: bounds,
				})
			}
		}
	}

	return result
}

// circleIntersectsRectangleは円が長方形と交差するかどうかをチェックします
// 円と長方形の交差判定に最適化されたアルゴリズムを使用します
func (g *RegionGridInfo) circleIntersectsRectangle(center Point, radiusSquared float64, rectMinX, rectMinY, rectMaxX, rectMaxY float64) bool {
	// 長方形上で円の中心に最も近い点を見つけます
	closestX := math.Max(rectMinX, math.Min(center.X, rectMaxX))
	closestY := math.Max(rectMinY, math.Min(center.Y, rectMaxY))

	// 円の中心から最も近い点までの距離の2乗を計算します
	dx := center.X - closestX
	dy := center.Y - closestY
	distanceSquared := dx*dx + dy*dy

	// 距離が半径の2乗以下であるかどうかを確認します
	return distanceSquared <= radiusSquared
}

// isCoordinateNearRegionは座標が領域の近くにあるかどうかをチェックします（拡大された境界内）
func (g *RegionGridInfo) isCoordinateNearRegion(coord LatLng, expandRadius float64) bool {
	// 拡大半径をメートルからおおよその度に変換します
	// これはパフォーマンスのための大まかな近似です
	approxDegreesPerMeter := 1.0 / 111000.0 // 赤道での大まかな近似値
	expandDegrees := expandRadius * approxDegreesPerMeter

	return coord.Lat <= g.TopLeft.Lat+expandDegrees &&
		coord.Lat >= g.BottomRight.Lat-expandDegrees &&
		coord.Lng >= g.TopLeft.Lng-expandDegrees &&
		coord.Lng <= g.BottomRight.Lng+expandDegrees
}

// GetGridsInRadiusWithDistanceは半径内のグリッドセルを中心からの距離とともに返します
func (g *RegionGridInfo) GetGridsInRadiusWithDistance(center LatLng, radiusMeters float64) ([]GridCellWithDistance, error) {
	cells, err := g.GetGridsInRadius(center, radiusMeters)
	if err != nil {
		return nil, err
	}

	result := make([]GridCellWithDistance, len(cells))
	centerMercator := g.latLngToWebMercator(center)

	for i, cell := range cells {
		// セルの中心までの距離を計算します
		cellCenterLat := (cell.Bounds.TopLeft.Lat + cell.Bounds.BottomRight.Lat) / 2
		cellCenterLng := (cell.Bounds.TopLeft.Lng + cell.Bounds.BottomRight.Lng) / 2
		cellCenterMercator := g.latLngToWebMercator(LatLng{Lat: cellCenterLat, Lng: cellCenterLng})

		dx := centerMercator.X - cellCenterMercator.X
		dy := centerMercator.Y - cellCenterMercator.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		result[i] = GridCellWithDistance{
			Cell:     cell,
			Distance: distance,
		}
	}

	return result, nil
}

// GridCellWithDistanceは参照点からの距離を持つグリッドセルを表します
type GridCellWithDistance struct {
	Cell     *GridCell
	Distance float64 // 距離（メートル）
}

// GetGridInfoはグリッドに関する基本情報を返します
func (g *RegionGridInfo) GetGridInfo() string {
	return fmt.Sprintf("グリッド: %dx%d セル, %.0fm セルサイズ, 地域: (%.6f,%.6f) から (%.6f,%.6f)",
		g.Cols, g.Rows, g.GridSizeMeters,
		g.TopLeft.Lat, g.TopLeft.Lng,
		g.BottomRight.Lat, g.BottomRight.Lng)
}

// Stringはグリッドセルの文字列表現を返します
func (gc *GridCell) String() string {
	return fmt.Sprintf("Grid[%d,%d] TopLeft(%.6f,%.6f) BottomRight(%.6f,%.6f)",
		gc.Row, gc.Col,
		gc.Bounds.TopLeft.Lat, gc.Bounds.TopLeft.Lng,
		gc.Bounds.BottomRight.Lat, gc.Bounds.BottomRight.Lng)
}

// 使用例
func main() {
	// 地域を定義：東京首都圏
	topLeft := LatLng{Lat: 35.8, Lng: 139.3}     // 北西の角
	bottomRight := LatLng{Lat: 35.6, Lng: 139.9} // 南東の角

	// 500mセルのグリッドを作成
	grid, err := NewRegionGridInfo(topLeft, bottomRight, 500.0)
	if err != nil {
		panic(err)
	}

	fmt.Println("地域グリッドシステムが作成されました")
	fmt.Println(grid.GetGridInfo())
	fmt.Println()

	// 地域内の座標をテスト
	testCoords := []LatLng{
		{Lat: 35.7, Lng: 139.4},   // 西部
		{Lat: 35.65, Lng: 139.5}, // 中央部
		{Lat: 35.75, Lng: 139.7}, // 東部
		{Lat: 35.68, Lng: 139.75},// 南東部
	}

	locations := []string{"西部", "中央部", "東部", "南東部"}

	for i, coord := range testCoords {
		fmt.Printf("=== %s ===\n", locations[i])

		// グリッドセルを取得
		cell, err := grid.GetGridCell(coord)
		if err != nil {
			fmt.Printf("エラー: %v\n", err)
			continue
		}

		fmt.Printf("入力: %.6f, %.6f\n", coord.Lat, coord.Lng)
		fmt.Printf("グリッド: %s\n", cell)

		// 高速版をテスト
		fastCell := grid.GetGridCellFast(coord)
		fmt.Printf("高速版: Grid[%d,%d]\n", fastCell.Row, fastCell.Col)

		// 隣接グリッドを取得
		neighbors := grid.GetNeighboringGrids(cell.Row, cell.Col)
		fmt.Printf("隣接グリッド数: %d\n", len(neighbors))

		fmt.Println()
	}

	// 地域外の座標をテスト
	fmt.Println("=== 地域外テスト ===")
	outsideCoord := LatLng{Lat: 36.0, Lng: 140.0} // 地域外
	_, err = grid.GetGridCell(outsideCoord)
	if err != nil {
		fmt.Printf("地域外座標に対する期待されるエラー: %v\n", err)
	}
	fmt.Println()

	// 半径検索のテスト
	fmt.Println("=== 半径検索テスト ===")
	searchCenter := LatLng{Lat: 35.7, Lng: 139.5} // 中心地
	searchRadius := 1000.0                       // 半径1km

	fmt.Printf("(%.6f, %.6f) から %.0fm 以内のグリッドを検索中\n", searchCenter.Lat, searchCenter.Lng, searchRadius)

	gridsInRadius, err := grid.GetGridsInRadius(searchCenter, searchRadius)
	if err != nil {
		fmt.Printf("半径検索でエラー: %v\n", err)
	} else {
		fmt.Printf("半径%.0fm以内に %d 個のグリッドが見つかりました\n", searchRadius, len(gridsInRadius))

		// 最初の数件の結果を表示
		maxShow := 5
		if len(gridsInRadius) > maxShow {
			fmt.Printf("最初の %d 個のグリッドを表示:\n", maxShow)
			for i := 0; i < maxShow; i++ {
				fmt.Printf("  %s\n", gridsInRadius[i])
			}
			fmt.Printf("  ... 他%d件\n", len(gridsInRadius)-maxShow)
		} else {
			for _, cell := range gridsInRadius {
				fmt.Printf("  %s\n", cell)
			}
		}
	}
	fmt.Println()

	// 距離付き半径検索のテスト
	fmt.Println("=== 距離付き半径検索テスト ===")
	gridsWithDistance, err := grid.GetGridsInRadiusWithDistance(searchCenter, searchRadius)
	if err != nil {
		fmt.Printf("距離付き半径検索でエラー: %v\n", err)
	} else {
		fmt.Printf("距離付きのグリッドが %d 個見つかりました:\n", len(gridsWithDistance))
		maxShow := 3
		for i := 0; i < len(gridsWithDistance) && i < maxShow; i++ {
			cell := gridsWithDistance[i]
			fmt.Printf("  %s, 距離: %.1fm\n", cell.Cell, cell.Distance)
		}
		if len(gridsWithDistance) > maxShow {
			fmt.Printf("  ... 他%d件\n", len(gridsWithDistance)-maxShow)
		}
	}
	fmt.Println()

	// 様々な半径サイズをテスト
	fmt.Println("=== 様々な半径サイズでのテスト ===")
	testRadii := []float64{500, 1000, 2000, 5000} // 500m, 1km, 2km, 5km
	for _, radius := range testRadii {
		grids, err := grid.GetGridsInRadius(searchCenter, radius)
		if err != nil {
			fmt.Printf("半径 %.0fm: エラー - %v\n", radius, err)
		} else {
			fmt.Printf("半径 %.0fm: %d グリッド\n", radius, len(grids))
		}
	}
	fmt.Println()

	// 高速版のテスト
	fmt.Println("=== 高速版テスト ===")
	fastGrids := grid.GetGridsInRadiusFast(searchCenter, searchRadius)
	fmt.Printf("高速版で見つかったグリッド数: %d (通常版と一致するはず: %d)\n", len(fastGrids), len(gridsInRadius))

	// グリッドセルの例をいくつか表示
	fmt.Println("=== グリッドセルサンプル ===")
	sampleCells := []*GridCell{}

	// 角のセルを取得
	if cell, err := grid.GetGridAtCoordinates(0, 0); err == nil {
		sampleCells = append(sampleCells, cell)
	}
	if cell, err := grid.GetGridAtCoordinates(0, grid.Cols-1); err == nil {
		sampleCells = append(sampleCells, cell)
	}
	if cell, err := grid.GetGridAtCoordinates(grid.Rows-1, 0); err == nil {
		sampleCells = append(sampleCells, cell)
	}
	if cell, err := grid.GetGridAtCoordinates(grid.Rows-1, grid.Cols-1); err == nil {
		sampleCells = append(sampleCells, cell)
	}

	corners := []string{"左上", "右上", "左下", "右下"}
	for i, cell := range sampleCells {
		fmt.Printf("%s: %s\n", corners[i], cell)
	}
}

// BenchmarkRegionGrid - パフォーマンステスト
func BenchmarkRegionGrid(coords []LatLng, grid *RegionGridInfo) {
	for i := 0; i < 1000000; i++ {
		for _, coord := range coords {
			grid.GetGridCellFast(coord)
		}
	}
}
