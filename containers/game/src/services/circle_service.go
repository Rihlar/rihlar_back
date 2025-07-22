package services

import (
	"fmt"
	"game/logger"
	"game/models"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type CircleService struct{}

//　円の詳細取得
func(CircleService) GetCircleDeteile(circleId string) (models.Circle, error) {

	//　円の詳細取得
	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
		return models.Circle{}, err
	}

	return circleDeteile, nil
}

// 画像取得
func(CircleService) GetCircleImage(circleId string) (string, error) {

	// circleIdから円のimageIdをとってくる
	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
		return "", err
	}

	// imageIdから画像パス生成　TODO:
	imagePath := "./assets/circle-images/" + circleDeteile.ImageID + ".png"

	return imagePath, nil
}
// 画像のアップロード
func (CircleService) UploadImage(circleId string, userId string, fileHeader *multipart.FileHeader) error {

	// 円取得
	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
		return err
	}

	// 円が自分が作成したかの判定
	if circleDeteile.UserID != userId {
		// 円の所有者じゃない場合
		err := fmt.Errorf("circle does not belong to user")
		logger.PrintErr("circle does not belong to user", err)
		return err
	}

	// フォルダパス
	uploadDir := "./assets/circle-images/"

	// 保存先ファイルパスの組み立て
	dstPath := filepath.Join(uploadDir, circleDeteile.ImageID+".png")

	// 画像がすでにあるのか確認
	if _, err := os.Stat(dstPath); err == nil {
	err := fmt.Errorf("file already exists: %s", dstPath)
	logger.PrintErr("file already exists", err)
	return err
}


	// アップロードされたファイルを開く
	src, err := fileHeader.Open()
	if err != nil {
		logger.PrintErr("ファイルオープン失敗", err)
		return err
	}
	defer src.Close()

	// ファイル保存(なんか外見だけを作って)
	dst, err := os.Create(dstPath)
	if err != nil {
		logger.PrintErr("File upload failed")
		return err
	}

	// 中身をコピー(こっちで中身を入れているイメージ。あってるかは知らん)
	written, err := io.Copy(dst, src)
	if err != nil {
		logger.PrintErr("コピー失敗", err)
		return err
	}
	// 書き込まれたか
	if written == 0 {
		return err
	}

	// ファイルのクローズ
	dst.Close()

	return nil
}

// 画像のリスト
func (CircleService) GetImageList(userid string) ([]string, error) {
	// プロフィールを取得
	profile,err := models.GetProfile(userid)
	if err != nil {
		logger.PrintErr("profile does not exist", err)
		return []string{}, err
	}

	// システムのゲームを取得
	sysGame,err := models.GetGameByID(profile.SysGame)

	// エラー処理
	if err != nil {
		logger.PrintErr("game does not exist", err)
		return []string{}, err
	}

	// ゲームの円を取得
	circles, err := sysGame.GetCircles()
	if err != nil {
		logger.PrintErr("circle does not exist", err)
		return []string{}, err
	}

	// 円のIDを格納する	
	circleIds := []string{}

	// 円を回す
	for _, circle := range circles {
		circleIds = append(circleIds, circle.CircleID)
	}

	return circleIds, nil
}