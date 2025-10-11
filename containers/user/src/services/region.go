package services

import "user/models"

// リージョン一覧を取得する
func GetAllRegions() ([]models.Region, error) {
	return models.GetRegions()
}
