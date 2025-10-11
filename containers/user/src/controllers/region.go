package controllers

import (
	"net/http"
	"user/services"

	"github.com/labstack/echo/v4"
)

func GetRegionList(ctx echo.Context) error {
	// サービスから地域情報を取得
	regionList, err := services.GetAllRegions()

	// エラー処理
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get region list"})
	}

	return ctx.JSON(http.StatusOK, regionList)
}
