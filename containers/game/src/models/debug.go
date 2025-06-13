package models

import (
	"game/logger"
)

// テーブルを削除してマイグレートする関数
func DeleteAndMigrate(table interface{}) error {
	// テーブルを削除する
	err := dbconn.Migrator().DropTable(table)

	// エラー処理
	if err != nil {
		logger.PrintErr("テーブルの削除中にエラーが発生しました",err)
	}

	// マイグレーションする
	return dbconn.AutoMigrate(table)
}

func Debug() {
	// マイグレーション のコードをここに書く
	DeleteAndMigrate(Sample{})

	// デバッグのコードを呼び出す
	DebugSample()
}
