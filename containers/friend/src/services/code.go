package services

import (
	"encoding/base64"
	"friend/logger"
	"friend/models"
	"friend/utils"

	"gorm.io/gorm"
)

// コードを生成する関数
func GenCode(userid string) (string,error) {
	// コードのもとになるUUID を生成する
	uid,err := utils.Genid()

	// エラー処理
	if err != nil {
		return "", err
	}

	// base58にエンコードする
	encoded, err := Encode([]byte(uid))

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return "", err
	}

	// データベースに保存する
	err = models.SaveCode(userid, encoded)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return "", err
	}

	return encoded, nil
}

// base64にエンコードする関数
func Encode(data []byte) (string,error) {
	return base64.URLEncoding.EncodeToString(data), nil
}

type FriendCode struct {
	Code     string 
	UseCount int    
}
// 現在のコードを取得する関数
func NowCode(userid string) (FriendCode,error) {
	// データベースから取得
	data,err := models.GetCode(userid)

	// コードが存在しないとき生成する
	if err == gorm.ErrRecordNotFound {
		// コードを生成する
		newCode, err := GenCode(userid)

		// エラー処理
		if err != nil {
			return FriendCode{}, err
		}

		return FriendCode{Code: newCode, UseCount: 0}, nil
	}

	// エラー処理
	if err != nil {
		return FriendCode{}, err
	}

	return FriendCode{Code: data.Code, UseCount: data.Count}, nil
}