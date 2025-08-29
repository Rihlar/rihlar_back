package models

// フレンド申請用のコードを管理するテーブル
type FriendCode struct {
	UserId string	`gorm:"primaryKey"`
	Code   string	`gorm:"type:varchar(50);unique"`
	UseCount int	`gorm:"default:0"`			//コードの使用回数
}
