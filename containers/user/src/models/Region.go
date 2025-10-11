package models

// テーブル定義
type Region struct {
	RegionID   string  `gorm:"primaryKey" json:"RegionID"`   // リージョンID
	RegionName string `gorm:"not null" json:"regionName"`	// リージョン名
	StartLat   float64 `gorm:"not null" json:"startLat"`	// 開始緯度
	StartLon   float64 `gorm:"not null" json:"startLon"`	// 開始経度
	EndLat     float64 `gorm:"not null" json:"endLat"`		// 終了緯度
	EndLon     float64 `gorm:"not null" json:"endLon"`		// 終了経度
}

// テーブル名
func (Region) TableName() string {
	return "regions"
}

// TODO デバッグ用
func GetRegionByID(id string) (Region, error) {
	region := Region{}

	// 取得
	err := dbconn.Where(&Region{
		RegionID: id,
	}).First(&region).Error

	// エラー処理
	if err != nil {
		return Region{}, err
	}

	return region, nil
}

// リージョン一覧を取得する
func GetRegions() ([]Region, error) {
	// データベースから読み込み
	regions := []Region{}

	// 取得
	err := dbconn.Find(&regions).Error

	// エラー処理
	if err != nil {
		return nil, err
	}

	return regions, nil
}
