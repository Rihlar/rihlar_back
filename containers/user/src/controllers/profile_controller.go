package controllers

import (
	"errors"
	
	"net/http"
	"user/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// updateするための構造体
type Input struct {
	Name     string `json:"name"`           //ユーザ名
	RecordID string `json:"record_id"`      // 実績ID
	Comment  string `json:"comment"`        //コメント
	RegionID string `json:"region_id"`      // 地域ID
	SysGame  string `json:"system_game_id"` // システムゲームID
	AdmGame  string `json:"admin_game_id"`  // アドミンゲームID
}

// プロフィールをIDから取得
func GetProfileById(c echo.Context) error {
	userID := c.Param("user_id")
	profile, err := services.GetProfileService(userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, profile)
}

// プロフィールをIDから変更
func UpdateProfileById(c echo.Context) error {
	userID := c.Param("user_id")
	var req Input

	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	input := services.Input{
		Name:     req.Name,
		RecordID: req.RecordID,
		Comment:  req.Comment,
		RegionID: req.RegionID,
		SysGame:  req.SysGame,
		AdmGame:  req.AdmGame,
	}

	if err := services.UpdateProfileById(userID, input); err != nil {
		// log.Println("サービス層エラー",err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "profile updated"})
}
