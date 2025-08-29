package models

type FriendType string

const (
	FriendTypeRequest FriendType = "request" //フレンドリクエストを送っている状態
	FriendTypeData    FriendType = "data"    //フレンドになった状態
)

// フレンドテーブル
type Friend struct {
	SenderId   string     `gorm:"primaryKey"` //フレンドの送信元ID
	ReceiverId string     `gorm:"primaryKey"` //フレンドの受信元ID
	Type       FriendType //データの種類
	CreatedAt  int64      //作成時間
}

// フレンドのデータを取得する関数
func GetFriendData(userid1, userid2 string) (Friend, error) {
	// データを格納する変数
	data := Friend{}

	// データベースを検索する
	err := dbconn.Where(&Friend{SenderId: userid1, ReceiverId: userid2}).Or(&Friend{SenderId: userid2, ReceiverId: userid1}).First(&data).Error

	// エラー処理
	if err != nil {
		return Friend{}, err
	}


	return data, nil
}

func SaveFriend(data Friend) error {
	return dbconn.Save(&data).Error
}

// フレンドリストを取得する関数
func GetFriendList(userid string) ([]Friend, error) {
	// データを格納する変数
	data := []Friend{}

	// データベースを検索する (送信元IDまたは受信元IDかつフレンドの種類)
	err := dbconn.Where(&Friend{SenderId: userid, Type: FriendTypeData}).Or(&Friend{ReceiverId: userid, Type: FriendTypeData}).Find(&data).Error

	// エラー処理
	if err != nil {
		return []Friend{}, err
	}

	return data, nil
}