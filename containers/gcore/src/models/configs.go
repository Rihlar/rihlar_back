package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	// Admin ゲーム一覧
	AdminGameId1 = "gameid-413a287b-213c-414f-a287-c1397db8f9bf"

	// システムゲーム一覧
	SysGameId1  = "gameid-9fcb784b-04a8-49c3-9ed9-ca9588eb86a8" // UserID1 に割り当て
	SysGameId2  = "gameid-b385def4-5833-4153-b54e-bdb134ab9fc8"
	SysGameId3  = "gameid-f9e8d7c6-b5a4-3210-fedc-ba9876543210"
	SysGameId4  = "gameid-e8d7c6b5-a4f3-2109-8765-43210fedcba9"
	SysGameId5  = "gameid-d7c6b5a4-f3e2-1098-7654-3210fedcba98"
	SysGameId6  = "gameid-c6b5a4f3-e2d1-0987-6543-210fedcba987"
	SysGameId7  = "gameid-b5a4f3e2-d1c0-9876-5432-10fedcba9876"
	SysGameId8  = "gameid-a4f3e2d1-c0b9-8765-4321-0fedcba98765"
	SysGameId9  = "gameid-f3e2d1c0-b9a8-7654-3210-fedcba987654"
	SysGameId10 = "gameid-e2d1c0b9-a8f7-6543-2109-fedcba987653"
	SysGameId11 = "gameid-d1c0b9a8-f7e6-5432-1098-fedcba987652"
	SysGameId12 = "gameid-c0b9a8f7-e6d5-4321-0987-fedcba987651"

	// UserID
	UserID1 = "userid-79541130-3275-4b90-8677-01323045aca5" // SysGameId1 に割り当て
	UserID2 = "userid-f224c85f-2ac4-4f50-a9dc-af80a659b671"
	UserID3 = "userid-a1b2c3d4-e5f6-7890-1234-567890abcdef"
	UserID4 = "userid-b2c3d4e5-f6a7-8901-2345-67890abcdef0"
	UserID5 = "userid-c3d4e5f6-a7b8-9012-3456-7890abcdef01"
	UserID6 = "userid-d4e5f6a7-b8c9-0123-4567-890abcdef012"
	UserID7 = "userid-e5f6a7b8-c9d0-1234-5678-90abcdef0123"
	UserID8 = "userid-f6a7b8c9-d0e1-2345-6789-0abcdef01234"
	UserID9 = "userid-a7b8c9d0-e1f2-3456-7890-abcdef012345"
	UserID10 = "userid-b8c9d0e1-f2a3-4567-890a-bcdef0123456"
	UserID11 = "userid-c9d0e1f2-a3b4-5678-90ab-cdef01234567"
	UserID12 = "userid-d0e1f2a3-b4c5-6789-0abc-def012345678"

	RegionId = "regionId-c161edb9-6aff-4244-8749-707bff2fa3be"
	GridMeter = 100

	// 緯度経度生成のための範囲
	MinLat = 33.5 // 南限
	MaxLat = 35.8 // 北限
	MinLon = 134.0 // 西限
	MaxLon = 137.5 // 東限

	RecordsPerUser = 100 // 各ユーザーごとの行動履歴件数 (基本件数)
	RandomCircleRecordsPerUser = 100 // 各ユーザーごとのランダムな円データ件数
	DuplicateLatLngRatio = 0.2 // 約20%の履歴で緯度経度を重複させる

	WalkRecordsPerUser = 50 // 各ユーザーごとの歩数記録件数

	// CSVファイルパス
	CircleDataCSVPath = "circle_data.csv"
)

// AllUserNames は定義されている全てのユーザー名のリストです。
var AllUserNames = []string{
	"Alice", "Bob", "Charlie", "David", "Eve", "Frank",
	"Grace", "Heidi", "Ivan", "Judy", "Kevin", "Liam",
}

// UserIDs は定義されている全てのユーザーIDのリストです。
var UserIDs = []string{
	UserID1, UserID2, UserID3, UserID4, UserID5, UserID6,
	UserID7, UserID8, UserID9, UserID10, UserID11, UserID12,
}

// SysGameIDs は定義されている全てのシステムゲームIDのリストです。
var SysGameIDs = []string{
	SysGameId1, SysGameId2, SysGameId3, SysGameId4, SysGameId5, SysGameId6,
	SysGameId7, SysGameId8, SysGameId9, SysGameId10, SysGameId11, SysGameId12,
}

// UserNameToIDMap はユーザー名と対応するユーザーIDのマップです。(参考用)
var UserNameToIDMap = map[string]string{
	"Alice":   UserID1,
	"Bob":     UserID2,
	"Charlie": UserID3,
	"David":   UserID4,
	"Eve":     UserID5,
	"Frank":   UserID6,
	"Grace":   UserID7,
	"Heidi":   UserID8,
	"Ivan":    UserID9,
	"Judy":    UserID10,
	"Kevin":   UserID11,
	"Liam":    UserID12,
}


// UserPointsMap は各ユーザーIDに紐づくポイントのマップです。
var UserPointsMap = map[string]int{
	UserID1: 1500, UserID2: 2300, UserID3: 800, UserID4: 3100,
	UserID5: 1200, UserID6: 2800, UserID7: 950, UserID8: 1900,
	UserID9: 600, UserID10: 2500, UserID11: 1750, UserID12: 1050,
}

// AllUserPoints は定義されている全てのユーザーのポイントのリストです。
var AllUserPoints = []int{
	1500, 2300, 800, 3100, 1200, 2800,
	950, 1900, 600, 2500, 1750, 1050,
}

// GlobalActionHistoryEntry は行動履歴の単一エントリを表す構造体です
type GlobalActionHistoryEntry struct {
	UserID    string
	Latitude  float64
	Longitude float64
	Steps     int
	Timestamp time.Time
}

// WalkRecord はユーザーの歩数記録を表す構造体です。
type WalkRecord struct {
	UserID    string
	Steps     int
	Timestamp time.Time
}

// AllUsersActionHistory は全てのユーザーの基本的な行動履歴をまとめた単一のリストです
var AllUsersActionHistory []GlobalActionHistoryEntry

// CircleDatas は各ユーザーの追加の行動履歴 (特定の場所データ) をまとめたリストです。
var CircleDatas []GlobalActionHistoryEntry

// AllWalkRecords は全てのユーザーの歩数記録をまとめた単一のリストです。
var AllWalkRecords []WalkRecord

// (以前のProvidedCircleDataは削除しました。代わりにCSVファイルを読み込みます)

var (
	loc *time.Location
	now time.Time
	startTime time.Time
)

func init() {
	rand.Seed(time.Now().UnixNano())

	loc, _ = time.LoadLocation("Asia/Tokyo")
	now = time.Now().In(loc)
	startTime = now.AddDate(0, -1, 0) // タイムスタンプ生成の基準となる日時 (例: 現在から1ヶ月前)

	// 各データ生成関数を呼び出す
	generateRandomActionHistory()
	generateRandomCircleDatas() // 各ユーザーにランダムな円データを追加
	generateCircleDataFromCSVFile() // CSVファイルから円データを追加
	generateRandomWalkRecords()

	// 生成された履歴の表示 (オプション)
	printGeneratedData()
}

// generateRandomActionHistory はランダムな行動履歴データを生成し、AllUsersActionHistoryに追加します。
func generateRandomActionHistory() {
	for _, userID := range UserIDs {
		sharedLocations := make(map[int]struct {
			Latitude  float64
			Longitude float64
		})

		// AllUsersActionHistory と CircleDatas (ランダム生成分とCSV読込分) の両方で重複を考慮するため、合計レコード数で計算
		// CSVデータ件数は generateCircleDataFromCSVFile の実行時に初めて確定するため、ここでは最大値として適当な値を設定するか、
		// またはCSV読み込み後に別途調整することも検討できます。
		// 簡単のため、ここではランダム生成分と平均的なCSV件数を仮定して重複計算を行います。
		estimatedCSVRowCount := 50 // CSVの行数を仮定
		totalPotentialRecords := RecordsPerUser + RandomCircleRecordsPerUser + estimatedCSVRowCount
		numSharedLocations := int(float64(totalPotentialRecords) * DuplicateLatLngRatio)
		if numSharedLocations == 0 && totalPotentialRecords > 0 {
			numSharedLocations = 1
		}
		for i := 0; i < numSharedLocations; i++ {
			lat := MinLat + rand.Float64()*(MaxLat-MinLat)
			lon := MinLon + rand.Float64()*(MaxLon-MinLon)
			sharedLocations[i] = struct {
				Latitude  float64
				Longitude float64
			}{Latitude: lat, Longitude: lon}
		}

		for i := 0; i < RecordsPerUser; i++ {
			var lat float64
			var lon float64

			if rand.Float64() < DuplicateLatLngRatio && len(sharedLocations) > 0 {
				keys := make([]int, 0, len(sharedLocations))
				for k := range sharedLocations {
					keys = append(keys, k)
				}
				selectedKey := keys[rand.Intn(len(keys))]
				loc := sharedLocations[selectedKey]
				lat = loc.Latitude
				lon = loc.Longitude
			} else {
				lat = MinLat + rand.Float64()*(MaxLat-MinLat)
				lon = MinLon + rand.Float64()*(MaxLon-MinLon)
			}

			steps := rand.Intn(901) + 100
			timestamp := startTime.Add(time.Duration(rand.Int63n(now.Sub(startTime).Nanoseconds())) * time.Nanosecond)

			AllUsersActionHistory = append(AllUsersActionHistory, GlobalActionHistoryEntry{
				UserID:    userID,
				Latitude:  lat,
				Longitude: lon,
				Steps:     steps,
				Timestamp: timestamp,
			})
		}
	}
}

// generateRandomWalkRecords はランダムな歩数記録データを生成し、AllWalkRecordsに追加します。
func generateRandomWalkRecords() {
	for _, userID := range UserIDs {
		for i := 0; i < WalkRecordsPerUser; i++ {
			steps := rand.Intn(10000) + 1000 // 1000歩から11000歩の範囲で生成
			timestamp := startTime.Add(time.Duration(rand.Int63n(now.Sub(startTime).Nanoseconds())) * time.Nanosecond)

			AllWalkRecords = append(AllWalkRecords, WalkRecord{
				UserID:    userID,
				Steps:     steps,
				Timestamp: timestamp,
			})
		}
	}
}

// generateRandomCircleDatas は各ユーザーにRandomCircleRecordsPerUser件のランダムな円データを生成します。
func generateRandomCircleDatas() {
	for _, userID := range UserIDs {
		sharedLocations := make(map[int]struct {
			Latitude  float64
			Longitude float64
		})

		// CircleDatas 専用の共有位置を生成
		numSharedLocations := int(float64(RandomCircleRecordsPerUser) * DuplicateLatLngRatio)
		if numSharedLocations == 0 && RandomCircleRecordsPerUser > 0 {
			numSharedLocations = 1
		}
		for i := 0; i < numSharedLocations; i++ {
			lat := MinLat + rand.Float64()*(MaxLat-MinLat)
			lon := MinLon + rand.Float64()*(MaxLon-MinLon)
			sharedLocations[i] = struct {
				Latitude  float64
				Longitude float64
			}{Latitude: lat, Longitude: lon}
		}

		for i := 0; i < RandomCircleRecordsPerUser; i++ {
			var lat float64
			var lon float64

			if rand.Float64() < DuplicateLatLngRatio && len(sharedLocations) > 0 {
				keys := make([]int, 0, len(sharedLocations))
				for k := range sharedLocations {
					keys = append(keys, k)
				}
				selectedKey := keys[rand.Intn(len(keys))]
				loc := sharedLocations[selectedKey]
				lat = loc.Latitude
				lon = loc.Longitude
			} else {
				lat = MinLat + rand.Float64()*(MaxLat-MinLat)
				lon = MinLon + rand.Float64()*(MaxLon-MinLon)
			}

			steps := rand.Intn(901) + 100
			timestamp := startTime.Add(time.Duration(rand.Int63n(now.Sub(startTime).Nanoseconds())) * time.Nanosecond)

			CircleDatas = append(CircleDatas, GlobalActionHistoryEntry{
				UserID:    userID,
				Latitude:  lat,
				Longitude: lon,
				Steps:     steps,
				Timestamp: timestamp,
			})
		}
	}
}

// generateCircleDataFromCSVFile は指定されたCSVファイルから円データをパースし、CircleDatasに追加します。
// 各データ行はユーザーIDを循環させて割り当てられます。
func generateCircleDataFromCSVFile() {
	f, err := os.Open(CircleDataCSVPath)
	if err != nil {
		fmt.Printf("CSVファイル '%s' の読み込みに失敗しました: %v\n", CircleDataCSVPath, err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4 // 緯度,経度,歩数,タイムスタンプ の4フィールドを期待

	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("CSVファイルの読み込みエラー: %v\n", err)
			continue
		}

		if len(record) != 4 {
			fmt.Printf("不正なCSV行のフォーマットです (期待: 4フィールド, 実際: %dフィールド): %v\n", len(record), record)
			continue
		}

		userID := UserIDs[i%len(UserIDs)] // ユーザーIDを循環させる

		lat, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			fmt.Printf("緯度のパースエラー ('%s'): %v\n", record[0], err)
			continue
		}
		lon, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			fmt.Printf("経度のパースエラー ('%s'): %v\n", record[1], err)
			continue
		}
		steps, err := strconv.Atoi(record[2])
		if err != nil {
			fmt.Printf("歩数のパースエラー ('%s'): %v\n", record[2], err)
			continue
		}
		unixTimestamp, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			fmt.Printf("タイムスタンプのパースエラー ('%s'): %v\n", record[3], err)
			continue
		}
		timestamp := time.Unix(unixTimestamp, 0).In(loc)

		CircleDatas = append(CircleDatas, GlobalActionHistoryEntry{
			UserID:    userID,
			Latitude:  lat,
			Longitude: lon,
			Steps:     steps,
			Timestamp: timestamp,
		})
		i++
	}
	fmt.Printf("CSVファイル '%s' から %d 件の円データを読み込みました。\n", CircleDataCSVPath, i)
}

// printGeneratedData は生成された各データセットの例をコンソールに出力します。
func printGeneratedData() {
	fmt.Println("--- 全ユーザーの基本行動履歴の例 (最初の10件) ---")
	for i, entry := range AllUsersActionHistory {
		if i >= 10 {
			break
		}
		fmt.Printf("  履歴 %d: UserID: %s, 緯度 %.4f, 経度 %.4f, 歩数 %d, タイムスタンプ %s\n",
			i+1, entry.UserID, entry.Latitude, entry.Longitude, entry.Steps, entry.Timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("\n全基本行動履歴件数: %d件 (1ユーザーあたり %d 件)\n",
		len(AllUsersActionHistory), RecordsPerUser)

	fmt.Println("\n--- 全ユーザーの追加行動履歴 (CircleDatas) の例 (最初の50件) ---")
	// CircleDatas はランダムデータとCSVデータの合計件数になる
	displayCount := len(CircleDatas)
	if displayCount > 50 {
		displayCount = 50
	}
	for i := 0; i < displayCount; i++ {
		entry := CircleDatas[i]
		fmt.Printf("  記録 %d: UserID: %s, 緯度 %.4f, 経度 %.4f, 歩数 %d, タイムスタンプ %s\n",
			i+1, entry.UserID, entry.Latitude, entry.Longitude, entry.Steps, entry.Timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("\n全追加行動履歴 (CircleDatas) 件数: %d件 (ランダム: 1ユーザーあたり %d件, CSV: %d件)\n",
		len(CircleDatas), RandomCircleRecordsPerUser, len(CircleDatas) - len(UserIDs) * RandomCircleRecordsPerUser)

	fmt.Println("\n--- 全ユーザーの歩数記録 (WalkRecords) の例 (最初の10件) ---")
	for i, record := range AllWalkRecords {
		if i >= 10 {
			break
		}
		fmt.Printf("  記録 %d: UserID: %s, 歩数 %d, タイムスタンプ %s\n",
			i+1, record.UserID, record.Steps, record.Timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("\n全歩数記録 (WalkRecords) 件数: %d件 (1ユーザーあたり %d 件)\n",
		len(AllWalkRecords), WalkRecordsPerUser)
}