package services

type CustomError struct {
	Code       int    //HTTPステータスコード
	LogMessage string //ログに出すメッセージ
	ErrMessage string //レスポンスに含めるエラーメッセージ
	Err        error  //実際のエラー
}
