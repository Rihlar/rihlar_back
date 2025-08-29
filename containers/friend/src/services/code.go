package services

import (
	"encoding/base64"
	"friend/logger"
	"friend/models"
	"friend/utils"
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