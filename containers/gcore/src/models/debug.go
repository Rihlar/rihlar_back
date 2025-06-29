package models

import (
	"gcore/logger"
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
	DeleteAndMigrate(Region{})
	DeleteAndMigrate(Member{})
	DeleteAndMigrate(Team{})
	DeleteAndMigrate(Game{})
	// DeleteAndMigrate(BaseChunk{})
	DeleteAndMigrate(GameChunk{})
	DeleteAndMigrate(Circle{})
	DeleteAndMigrate(MovementLog{})
	DeleteAndMigrate(Profile{})

	// デバッグのコードを呼び出す
	DebugSample()

	DebugProfile()
	// DebugTeam()
	// DebugMember()

	// ゲームで使用するリージョンを作成する
	DebugRegion()

	// ゲームをデバッグする
	DebugGame()

	// DebugBaseChunk()
	DebugGameChunk()

	DebugPerformance()
	// DebugCircle()
	// DebugMovementLog()
}

