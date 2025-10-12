package controllers

import (
	"errors"

	"net/http"
	"user/logger"
	"user/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

/*
	プロフィール処理
*/

// 管理者がプロフィールを作成する用の構造体
type DebugAdminProfile struct {
	UserID string	`json:"user_id"`
	services.Input	// 基本プロフィール
}

// プロフィール作成
func CreateProfileFromAdmin(ctx echo.Context) error {
	var req DebugAdminProfile
	//リクエスト整形(bind)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//作成
	userID, err := services.CreateProfileService(req.UserID, req.Input)

	//エラー処理
	if err != nil{
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"status": "failed to create"})
	}

	//成功時ユーザーIDを返す
	return ctx.JSON(http.StatusCreated, echo.Map{
		"status": "profile created",
		"user_id": userID,
	})
}

// プロフィールをIDから変更 (管理者操作用)
func UpdateProfileByIdFromAdmin(ctx echo.Context) error {
	//useridをヘッダーから取得
	// userID := c.Request().Header.Get("userid")
	// userID := ctx.Get("UserID").(string)
	var req DebugAdminProfile

	//リクエスト整形(bind)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdateProfileById(req.UserID, req.Input); err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		
		if err.Error() == "no rows updated" {
			return ctx.JSON(http.StatusConflict, echo.Map{"error": "no changes detected"})
		}
		//NotFound以外のエラー
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"status": "profile updated"})
}

// プロフィール一覧を取得
func GetAllProfilesFromAdmin(ctx echo.Context) error {
	//エラー処理
	profiles, err := services.GetAllProfilesService()
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"status": "failed to get profiles"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"profiles": profiles,
	})
}

// プロファイルを削除する
func DeleteProfileFromAdmin(ctx echo.Context) error {
	// ヘッダから削除対象のユーザーIDを取得
	userID := ctx.Request().Header.Get("UserID")

	//エラー処理
	if err := services.DeleteProfileService(userID); err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		//NotFound以外のエラー
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "delete failed"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"status": "profile deleted"})
}

func GetProfileFromAdmin(c echo.Context) error {
	//useridをヘッダーから取得
	userID := c.Request().Header.Get("UserID")

	//プロフィール取得
	profile, err := services.GetProfileService(userID)

	//エラー処理
	if err != nil {
		//結果がNotFoundの時
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		logger.PrintErr(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, profile)
}