package models

import "gcore/logger"

var (
	// TODO とりあえず現状のリージョンをかく
	regions = []Region{
		{
			RegionID:   "regionId-a3d3dcd1-7a73-4a31-9908-b9cab944280d",
			RegionName: "北海道",
			StartLat:   45.55,
			StartLon:   139.5,
			EndLat:     41.3,
			EndLon:     148.5,
		},
		{
			RegionID:   "regionId-65051f6a-9e94-439d-a7f6-5c127ad0c885",
			RegionName: "東北",
			StartLat:   41.6,
			StartLon:   138.0,
			EndLat:     37.0,
			EndLon:     142.5,
		},
		{
			RegionID:   "regionId-fb145c05-e0e5-4f22-86e1-9f40326faf31",
			RegionName: "関東",
			StartLat:   37.0,
			StartLon:   138.5,
			EndLat:     34.5,
			EndLon:     141.0,
		},
		{
			RegionID:   "regionId-16d687a3-8eab-4b36-8563-b62514823fe8",
			RegionName: "中部",
			StartLat:   38.0,
			StartLon:   135.5,
			EndLat:     34.0,
			EndLon:     139.5,
		},
		{
			RegionID:   RegionId,
			RegionName: "関西",
			StartLat:   35.8,
			StartLon:   134.0,
			EndLat:     33.5,
			EndLon:     137.5,
		},
		{
			RegionID:   "regionId-005344a4-523d-4aab-9d1e-d232322cf54e",
			RegionName: "中国",
			StartLat:   36.0,
			StartLon:   130.5,
			EndLat:     33.5,
			EndLon:     134.5,
		},
		{
			RegionID:   "regionId-3501902d-8bab-40cd-926f-30c53d80efc5",
			RegionName: "四国",
			StartLat:   34.5,
			StartLon:   132.0,
			EndLat:     32.5,
			EndLon:     135.5,
		},
		{
			RegionID:   "regionId-ef5aa179-53e0-481d-b64d-ae7654049a88",
			RegionName: "九州",
			StartLat:   34.5,
			StartLon:   128.0,
			EndLat:     30.5,
			EndLon:     132.5,
		},
	}
)


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

func InitRegion() error {
	// リージョンを登録する
	for _, region := range regions {
		// 書き込み
		err := dbconn.Save(&region).Error

		// エラー処理
		if err != nil {
			return err
		}
	}

	return nil
}

func DebugRegion() {
	// デバッグ用のコードをここに書く
	// 書き込み
	InitRegion()

	logger.Println("リージョン取得成功")
}
