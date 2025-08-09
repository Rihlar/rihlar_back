package services

import (
	"gcore/logger"
	"gcore/models"
	"net/http"
	"time"
)

// 歩いたことを報告する引数
type MovementArgs struct {
	UserID    string  `json:"userID`    // ユーザーID
	Steps     int64   `json:"steps`     // 歩数
	Latitude  float64 `json:"latitude`  // 緯度
	Longitude float64 `json:"longitude` // 経度
	TimeStamp int64   `json:"timeStamp`
}

type SaveMovementLogArgs struct {
	UserID    string        `json:"userID`    // ユーザーID
	SystemID  string        `json:"systemID`  // システムゲームID
	Games     []models.Game `json:"games`     // 管理ゲームID
	Steps     int64         `json:"steps`     // 歩数
	Latitude  float64       `json:"latitude`  // 緯度
	Longitude float64       `json:"longitude` // 経度
	TimeStamp int64         `json:"timeStamp`
}

type ProcessChunkArgs struct {
	UserID    string        `json:"userID`    // ユーザーID
	Latitude  float64       `json:"latitude`  // 緯度
	Longitude float64       `json:"longitude` // 経度
	Games     []models.Game `json:"games`     // 管理ゲームID
}

func ReportMovement(args MovementArgs) (ProcessChunkResponse,error) {
	// プロファイルを取得する
	profile, err := models.GetProfile(args.UserID)

	// エラー処理
	if err != nil {
		return ProcessChunkResponse{}, err
	}

	// システムゲームを取得
	sysGame, err := models.GetGame(profile.SysGame)

	// エラー処理
	if err != nil {
		return ProcessChunkResponse{}, err
	}

	// adminゲームを取得
	admGame, err := models.GetGame(profile.AdmGame)

	// エラー処理
	if err != nil {
		return ProcessChunkResponse{}, err
	}

	// タイムスタンプがなかったら
	if args.TimeStamp == 0 {
		args.TimeStamp = time.Now().Unix()
	}

	// 行動記録を保存する
	err = SaveMovementLog(SaveMovementLogArgs{
		UserID:    args.UserID,
		SystemID:  sysGame.GameID,
		Games:     []models.Game{admGame, sysGame},
		Steps:     args.Steps,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		TimeStamp: args.TimeStamp,
	})

	// エラー処理
	if err != nil {
		return ProcessChunkResponse{}, err
	}

	// チャンクに対しての処理
	response,err := ProcessChunk(ProcessChunkArgs{
		UserID:    args.UserID,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     []models.Game{admGame, sysGame},
	})

	// エラー処理
	if err != nil {
		return ProcessChunkResponse{}, err
	}

	// 円のレベルアップ処理
	err = LevelUpCircle(ProcessCircleArgs{
		UserID:    args.UserID,
		Latitude:  args.Latitude,
		Longitude: args.Longitude,
		Games:     []models.Game{admGame, sysGame},
		Steps:     args.Steps,
	})

	// エラー処理
	if err != nil {
		return ProcessChunkResponse{}, err
	}

	return response, nil
}

// 歩いたログを記録する関数
func SaveMovementLog(args SaveMovementLogArgs) error {
	// ゲームを回す
	for _, game := range args.Games {
		// メンバーオブジェクト取得
		member, err := game.GetMemberByUserID(args.UserID)

		// エラー処理
		if err != nil {
			return err
		}

		// 歩いたログを保存する
		if err := member.SaveMovementLog(args.Latitude, args.Longitude, args.Steps, args.TimeStamp); err != nil {
			logger.PrintErr(err)
			return err
		}
	}

	return nil
}

type ProcessChunkAdminGameResponse struct {
	IsSuccess bool   //管理ゲームが成功したか
	GameId    string //管理ゲームID
	Message   string //メッセージ
	Status    int    //ステータス
}

type ProcessChunkResponse struct {
	IsSyetemSuccess bool //システムゲームが成功したか
	AdminGames      []ProcessChunkAdminGameResponse	//管理ゲームの処理結果
}

// チャンクに対しての処理 (Level1 などを実行する)
func ProcessChunk(args ProcessChunkArgs) (ProcessChunkResponse,error) {
	returnData := ProcessChunkResponse{}

	// ゲームを回す
	for _, game := range args.Games {
		// 一番近いチャンクを取得
		chunk, err := game.GetChunkByLatLon(args.Latitude, args.Longitude)

		// エラー処理
		if err != nil {
			// エラーが起きた時
			// admin ゲームならむし
			if game.Type == 1 {
				logger.Println("チャンク処理でエラーが発生しました")
				logger.Println(err)

				// admin ゲームの追加するデータ
				addData := ProcessChunkAdminGameResponse{
					IsSuccess: false,
					Message:   err.Error(),
					Status:    http.StatusInternalServerError,
					GameId:    game.GameID,
				}

				// 追加する
				returnData.AdminGames = append(returnData.AdminGames, addData)
			} else {
				// システムゲームの時
				// レスポンスを変更する
				returnData.IsSyetemSuccess = false

				// エラーを返す
				return returnData, err
			}

			continue
		}

		logger.Println("near chunk:", chunk)

		// チャンクのレベルが0か1なら
		if chunk.Level == 0 || chunk.Level == 1 {
			// チャンクを更新する
			if err := chunk.ChangeLevel(1); err != nil {
				logger.PrintErr(err)
				return returnData,err
			}

			// オーナーも変更する
			if err := chunk.ChangeOwner(args.UserID); err != nil {
				logger.PrintErr(err)
				return returnData,err
			}
		}
	}

	return returnData, nil
}

type MovementLog struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Steps     int64   `json:"steps"`
	TimeStamp int64   `json:"timeStamp"`
}

// 歩いた記録を取得する
func GetReportedMovement(gameId, userId string) ([]MovementLog, error) {
	// ゲームを取得する
	game, err := models.GetGame(gameId)

	// エラー処理
	if err != nil {
		return []MovementLog{}, err
	}

	// メンバーオブジェクト取得
	member, err := game.GetMemberByUserID(userId)

	// エラー処理
	if err != nil {
		return []MovementLog{}, err
	}

	// 歩いたログを保存する
	movementLogs, err := member.GetReportedMovement()

	// エラー処理
	if err != nil {
		return []MovementLog{}, err
	}

	returnDatas := []MovementLog{}

	for _, movementLog := range movementLogs {
		// 返却用に変換
		returnDatas = append(returnDatas, MovementLog{
			Latitude:  movementLog.Latitude,
			Longitude: movementLog.Longitude,
			Steps:     movementLog.Steps,
			TimeStamp: movementLog.TimeStamp,
		})
	}

	return returnDatas, nil
}
