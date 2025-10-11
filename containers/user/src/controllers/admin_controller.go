package controllers

import (
	"net/http"
	"user/models"

	"github.com/labstack/echo/v4"
)

// 全てのユーザー情報を取得する
func GetAllUsersFromAdmin(c echo.Context) error {
	// モデルから全てのプロフィールを取得
	profiles, err := models.GetAllProfiles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to get profiles"})
	}

	return c.JSON(http.StatusOK, profiles)
}
