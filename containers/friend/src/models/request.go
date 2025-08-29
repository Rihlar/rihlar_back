package models

// 受信済みリクエストを取得する関数
func GetRecvedRequests(userId string) ([]Friend, error) {
	// 結果を格納する関数
	datas := []Friend{}

	// ユーザーIDから検索する
	err := dbconn.Where(&Friend{
		ReceiverId: userId,
		Type:       FriendTypeRequest,
	}).Find(&datas).Error

	// エラー処理
	if err != nil {
		return []Friend{}, err
	}

	return datas, nil
}

// 送信済みリクエストを取得する関数
func GetSentRequests(userId string) ([]Friend, error) {
	// 結果を格納する関数
	datas := []Friend{}

	// ユーザーIDから検索する
	err := dbconn.Where(&Friend{
		SenderId:   userId,
		Type:       FriendTypeRequest,
	}).Find(&datas).Error

	// エラー処理
	if err != nil {
		return []Friend{}, err
	}

	return datas, nil
}