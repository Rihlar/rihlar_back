package services

import (
	"errors"
	"friend/logger"
	"friend/models"
)

type RequestData struct {
	UserId    string `json:"userId"`
	UserName  string `json:"userName"`
	TimeStamp int64  `json:"timeStamp"`
}

func GetRecvedRequest(userId string) ([]RequestData, error) {
	// リクエストを取得する
	datas, err := models.GetRecvedRequests(userId)

	// エラー処理
	if err != nil {
		return []RequestData{}, err
	}

	// 返すデータを格納する変数
	returnDatas := []RequestData{}

	// リクエストを回して取得する
	for _, data := range datas {
		// 送信者の情報を取得する
		senderProfile,err := models.GetProfile(data.SenderId)

		// エラー処理
		if err != nil {
			logger.PrintErr(err)
			continue
		}

		// 返すデータを格納する
		returnDatas = append(returnDatas, RequestData{
			UserId:    data.SenderId,
			UserName:  senderProfile.Name,
			TimeStamp: data.CreatedAt,
		})
	}

	return returnDatas, nil
}

// リクエストを拒否する関数
func RejectRequest(userId, targetUserId string) error {
	// フレンドデータを取得
	friendData, err := models.GetFriendData(userId, targetUserId)

	// エラー処理
	if err != nil {
		return err
	}

	// リクエストかどうか
	if friendData.Type != models.FriendTypeRequest {
		// リクエストじゃない場合
		return errors.New("invalid request")
	}

	// 自信が受信者かどうか
	if friendData.ReceiverId != userId {
		// 自身当てのリクエストじゃない場合
		return errors.New("invalid request")
	}

	// フレンドデータを削除
	return models.DeleteFriend(friendData)
}