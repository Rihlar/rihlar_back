package services

import (
	"errors"
	"friend/models"
	"time"

	"gorm.io/gorm"
)

const (
	CodeMaxUseCount = 5		//コードの最大使用回数
)

// フレンドリクエストを送信する関数
func SendRequest(userId string, friendCode string) error {
	// コードから送信先の情報を取得する
	targetData, err := models.GetFromCode(friendCode)

	// エラー処理
	if err != nil {
		return err
	}

	// 自信と同じユーザーに送っている場合
	if userId == targetData.UserId {
		return errors.New("cannot send to yourself")
	}

	// コードの使用回数を検証する
	if targetData.UseCount >= CodeMaxUseCount {
		return errors.New("code is expired")
	}

	// フレンドのデータを取得する
	data, err := models.GetFriendData(userId, targetData.UserId)

	// レコードが存在しない場合
	if err == gorm.ErrRecordNotFound {
		// 使用回数を1回増やす
		err = models.UpdateUseCount(targetData.UserId, targetData.UseCount+1)
		
		// エラー処理
		if err != nil {
			return err
		}

		// リクエストを送る
		return models.SaveFriend(models.Friend{
			SenderId:   userId,                   //送信者ID 自身のID
			ReceiverId: targetData.UserId,        //受信者ID コードから取得
			Type:       models.FriendTypeRequest, //リクエスト
			CreatedAt:  time.Now().Unix(),        //送信時間
		})
	}

	// エラー処理
	if err != nil {
		return err
	}

	// すでにリクエストを送っている場合
	if data.Type == models.FriendTypeRequest {
		return errors.New("already request")
	}

	// すでにフレンドの場合
	if data.Type == models.FriendTypeData {
		return errors.New("already friend")
	}

	// 使用回数を1回増やす
	err = models.UpdateUseCount(targetData.UserId, targetData.UseCount+1)
	
	// エラー処理
	if err != nil {
		return err
	}

	// リクエストを送る
	return models.SaveFriend(models.Friend{
		SenderId:   userId,                   //送信者ID 自身のID
		ReceiverId: targetData.UserId,        //受信者ID コードから取得
		Type:       models.FriendTypeRequest, //リクエスト
		CreatedAt:  time.Now().Unix(),        //送信時間
	})
}
