package models

import (
	"fmt"
	"math/rand"
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
	AdditionalRecordsPerUser = 10 // 各ユーザーごとに追加する行動履歴件数 (circleDatas用)
	DuplicateLatLngRatio = 0.2 // 約20%の履歴で緯度経度を重複させる
)

// AllUserNames は定義されている全てのユーザー名のリストです。
var AllUserNames = []string{
	"Alice",
	"Bob",
	"Charlie",
	"David",
	"Eve",
	"Frank",
	"Grace",
	"Heidi",
	"Ivan",
	"Judy",
	"Kevin",
	"Liam",
}

// UserIDs は定義されている全てのユーザーIDのリストです。
var UserIDs = []string{
	UserID1,
	UserID2,
	UserID3,
	UserID4,
	UserID5,
	UserID6,
	UserID7,
	UserID8,
	UserID9,
	UserID10,
	UserID11,
	UserID12,
}

// SysGameIDs は定義されている全てのシステムゲームIDのリストです。
var SysGameIDs = []string{
	SysGameId1,
	SysGameId2,
	SysGameId3,
	SysGameId4,
	SysGameId5,
	SysGameId6,
	SysGameId7,
	SysGameId8,
	SysGameId9,
	SysGameId10,
	SysGameId11,
	SysGameId12,
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
	UserID1: 1500, // Alice
	UserID2: 2300, // Bob
	UserID3: 800,  // Charlie
	UserID4: 3100, // David
	UserID5: 1200, // Eve
	UserID6: 2800, // Frank
	UserID7: 950,  // Grace
	UserID8: 1900, // Heidi
	UserID9: 600,  // Ivan
	UserID10: 2500, // Judy
	UserID11: 1750, // Kevin
	UserID12: 1050, // Liam
}

// AllUserPoints は定義されている全てのユーザーのポイントのリストです。
var AllUserPoints = []int{
	1500, // UserID1 (Alice)
	2300, // UserID2 (Bob)
	800,  // UserID3 (Charlie)
	3100, // UserID4 (David)
	1200, // UserID5 (Eve)
	2800, // UserID6 (Frank)
	950,  // UserID7 (Grace)
	1900, // UserID8 (Heidi)
	600,  // UserID9 (Ivan)
	2500, // UserID10 (Judy)
	1750, // UserID11 (Kevin)
	1050, // UserID12 (Liam)
}

// GlobalActionHistoryEntry は行動履歴の単一エントリを表す構造体です
type GlobalActionHistoryEntry struct {
	UserID    string
	Latitude  float64
	Longitude float64
	Steps     int
	Timestamp time.Time
}

// AllUsersActionHistory は全てのユーザーの基本的な行動履歴をまとめた単一のリストです
var AllUsersActionHistory []GlobalActionHistoryEntry

// CircleDatas は各ユーザーの追加の行動履歴 (10件) をまとめたリストです。
// 通常の行動履歴とは区別して使用されることを想定しています。
var CircleDatas []GlobalActionHistoryEntry

func init() {
	rand.Seed(time.Now().UnixNano())

	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)
	// タイムスタンプ生成の基準となる日時 (例: 現在から1ヶ月前)
	startTime := now.AddDate(0, -1, 0)

	for _, userID := range UserIDs {
		// このユーザーの履歴で、緯度経度を共有するためのマップ
		sharedLocations := make(map[int]struct {
			Latitude  float64
			Longitude float64
		})

		// 共有する緯度経度を事前にいくつか生成
		// 基本履歴と追加履歴の両方で共有できるように考慮
		totalRecordsForUser := RecordsPerUser + AdditionalRecordsPerUser
		numSharedLocations := int(float64(totalRecordsForUser) * DuplicateLatLngRatio)
		if numSharedLocations == 0 && totalRecordsForUser > 0 {
			numSharedLocations = 1 // 少なくとも1つは共有地点を作る
		}
		for i := 0; i < numSharedLocations; i++ {
			lat := MinLat + rand.Float64()*(MaxLat-MinLat)
			lon := MinLon + rand.Float64()*(MaxLon-MinLon)
			sharedLocations[i] = struct {
				Latitude  float64
				Longitude float64
			}{Latitude: lat, Longitude: lon}
		}

		// RecordsPerUser件の基本的な行動履歴を生成し、AllUsersActionHistoryに追加
		for i := 0; i < RecordsPerUser; i++ {
			var lat float64
			var lon float64

			// 一定の確率で既存の緯度経度を再利用
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

		// AdditionalRecordsPerUser件の追加行動履歴を生成し、CircleDatasに追加
		for i := 0; i < AdditionalRecordsPerUser; i++ {
			var lat float64
			var lon float64

			// こちらも一定の確率で既存の緯度経度を再利用
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

	// 生成された履歴の一部を表示して確認 (オプション)
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

	fmt.Println("\n--- 全ユーザーの追加行動履歴 (CircleDatas) の例 (最初の10件) ---")
	for i, entry := range CircleDatas {
		if i >= 10 {
			break
		}
		fmt.Printf("  履歴 %d: UserID: %s, 緯度 %.4f, 経度 %.4f, 歩数 %d, タイムスタンプ %s\n",
			i+1, entry.UserID, entry.Latitude, entry.Longitude, entry.Steps, entry.Timestamp.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("\n全追加行動履歴 (CircleDatas) 件数: %d件 (1ユーザーあたり %d 件)\n",
		len(CircleDatas), AdditionalRecordsPerUser)
}