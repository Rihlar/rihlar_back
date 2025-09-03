package models

// フレンド申請用のコードを管理するテーブル
type FriendCode struct {
	UserId   string `gorm:"primaryKey"`
	Code     string `gorm:"type:varchar(50);unique"`
	UseCount int    `gorm:"default:0"` //コードの使用回数
}

// コードを保存する関数
func SaveCode(userId, code string) error {
	// データを保存する
	return dbconn.Save(&FriendCode{UserId: userId, Code: code, UseCount: 0}).Error
}

// コードの使用回数を更新する関数
func UpdateUseCount(userId string, Count int) error {
	// データを更新する
	return dbconn.Where(&FriendCode{UserId: userId}).Updates(FriendCode{UseCount: Count}).Error
}

// ユーザーIDごとの使用回数を取得する関数
func GetUseCount(userId string) (int, error) {
	var code FriendCode

	result := dbconn.Where(&FriendCode{UserId: userId}).First(&code)

	if result.Error != nil {
		return 0, result.Error
	}

	return code.UseCount, nil
}

type CodeData struct {
	Code   string `json:"code"`    // コード
	Count  int    `json:"count"`   // 使用回数
}

func GetCode(userId string) (CodeData, error) {
	var code FriendCode

	// データベースから取得
	result := dbconn.Where(&FriendCode{UserId: userId}).First(&code)

	if result.Error != nil {
		return CodeData{}, result.Error
	}

	// コード情報を返す
	return CodeData{Code: code.Code, Count: code.UseCount}, nil
}

// コードから情報を取得する
func GetFromCode(code string) (FriendCode, error) {
	var friendCode FriendCode

	result := dbconn.Where(&FriendCode{Code: code}).First(&friendCode)

	if result.Error != nil {
		return FriendCode{}, result.Error
	}

	return friendCode, nil
}