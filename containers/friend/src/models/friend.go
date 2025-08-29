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
