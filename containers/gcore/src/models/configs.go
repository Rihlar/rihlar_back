package models

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