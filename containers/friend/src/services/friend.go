package services

import (
	"errors"
	"friend/logger"
	"friend/models"
	"time"

	"gorm.io/gorm"
)

const (
	CodeMaxUseCount = 5 //コードの最大使用回数
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

type FriendData struct {
	UserId    string `json:"userId"`    //ユーザーID
	UserName  string `json:"userName"`  //ユーザー名
	TimeStamp int64  `json:"timeStamp"` //フレンドのなった時間
}

// フレンドのリストを取得する関数
func GetFriendList(userId string) ([]FriendData, error) {
	// フレンドリストを取得する
	datas, err := models.GetFriendList(userId)

	// エラー処理
	if err != nil {
		return []FriendData{}, err
	}

	// 返すデータを格納する変数
	returnDatas := []FriendData{}

	for _, data := range datas {
		// 相手のID
		targetId := ""

		// 相手のIDを特定する
		if data.SenderId == userId {
			// 自信が送信者の場合相手が受信者になる
			targetId = data.ReceiverId
		} else {
			// 自信が受信者の場合相手が送信者になる
			targetId = data.SenderId
		}

		// 相手のプロファイルを取得
		profile, err := models.GetProfile(targetId)

		// エラー処理
		if err != nil {
			logger.PrintErr("フレンド情報取得エラー", err)
			continue
		}

		// 返すデータを格納
		returnDatas = append(returnDatas, FriendData{
			UserId:    targetId,       //ユーザーID
			UserName:  profile.Name,   //ユーザー名
			TimeStamp: data.CreatedAt, //フレンドのなった時間
		})
	}

	// 返すデータを返す
	return returnDatas, nil
}
