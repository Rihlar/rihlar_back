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
	//useridをヘッダーから取得
	userID := c.Request().Header.Get("userid")

	//プロフィール取得
	profile, err := services.GetProfileService(userID)

	//エラー処理
	if err != nil {
		//結果がNotFoundの時
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, profile)
}

func CreateProfile(c echo.Context) error {
	// ユーザーIDを取得する	
	// userID := c.Request().Header.Get("userid")
	UserID := c.Get("UserID").(string)

	var req Input
	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//入力したデータを格納
	input := services.Input{
		Name:     req.Name,
		RecordID: req.RecordID,
		Comment:  req.Comment,
		RegionID: req.RegionID,
		SysGame:  req.SysGame,
		AdmGame:  "",
	}

	//作成
	userID, err := services.CreateProfileService(UserID,input)

	//エラー処理
	if err != nil{
		return c.JSON(http.StatusInternalServerError, echo.Map{"status": "failed to create"})
	}

	//成功時ユーザーIDを返す
	return c.JSON(http.StatusCreated, echo.Map{
		"status": "profile created",
		"user_id": userID,
	})
}

// プロフィールをIDから変更
func UpdateProfileById(c echo.Context) error {
	//useridをヘッダーから取得
	userID := c.Request().Header.Get("userid")
	var req Input

	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//リクエストを格納
	input := services.Input{
		Name:     req.Name,
		RecordID: req.RecordID,
		Comment:  req.Comment,
		RegionID: req.RegionID,
	}

	//エラー処理
	if err := services.UpdateProfileById(userID, input); err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "profile updated"})
}

// プライバシー情報を取得
func GetPrivacyProfile(c echo.Context) error {
	//useridをヘッダーから取得
	UserID := c.Request().Header.Get("userid")

	//取得
	privacy, err := services.GetPrivacyProfileService(UserID)

	//エラー処理
	if err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "privacy profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, privacy)
}

// プライバシー情報を更新する
func UpdatePrivacyProfile(c echo.Context) error {

	//useridをヘッダーから取得
	UserID := c.Request().Header.Get("userid")
	var req services.PrivacyInput

	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdatePrivacyProfileById(UserID, req); err != nil {
		//Notfoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "privacy profile not found"})
		}
		//Notfound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "privacy updated!"})
}

// 地域情報の取得
func GetRegionProfile(c echo.Context) error {

	//useridをヘッダーから取得
	UserID := c.Request().Header.Get("userID")

	//取得
	region, err := services.GetRegionProfileService(UserID)

	//エラー処理
	if err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "region profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, region)
}

// 地域情報の編集
func UpdateRegionProfile(c echo.Context) error {
	//useridをヘッダーから取得
	UserID := c.Request().Header.Get("Userid")

	//格納用の地域情報
	var req struct {
		RegionID string `json:"region_id"`
	}
	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdateRegionById(UserID, req.RegionID); err != nil {
		//NotFOundのとき
		if errors.Is(err,gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "region profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}
	
	return c.JSON(http.StatusOK, echo.Map{"status": "region profile updated!"})
}
