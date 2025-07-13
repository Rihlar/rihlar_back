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

	// imageIdから画像パス生成
	imagePath := "./assets/circle-images/" + circleDeteile.ImageID + ".png"

	return imagePath, nil
}
// 画像のアップロード
func (CircleService) UploadImage(circleId string, fileHeader *multipart.FileHeader) error {

	// 円取得
	circleDeteile, err := models.GetCircleDeteile(circleId)
	if err != nil {
		logger.PrintErr("circle deteile does not exist", err)
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
