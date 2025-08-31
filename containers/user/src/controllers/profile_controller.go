package controllers

import (
	"errors"

	"net/http"
	"user/services"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)


/*
	プロフィール処理
*/

// プロフィールをIDから取得
func GetProfileById(c echo.Context) error {
	//useridをヘッダーから取得
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)

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

// プロフィール作成
func CreateProfile(c echo.Context) error {
	// ユーザーIDを取得する	
	// userID := c.Request().Header.Get("userid")
	UserID := c.Get("UserID").(string)

	var req services.Input
	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//作成
	userID, err := services.CreateProfileService(UserID, req)

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
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)
	var req services.Input

	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdateProfileById(userID, req); err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "profile updated"})
}

/*
	実績関連処理
*/

//実績を取得する
func GetAchiveProfile(c echo.Context) error {
	//userIDを取得
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)

	//service実行
	Achivement, err := services.GetAchiveProfile(userID)

	//エラー処理
	if err != nil {
		//NotFoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Achivement profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal error"})
	}
	
	//実績を返す
	return c.JSON(http.StatusOK, Achivement)
}

// 実績更新
func UpdateAchiveProfile(c echo.Context) error {
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)
	var req services.AchiveInput

	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdateAchiveProfile(userID, req); err != nil {
		//Notfoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Achivement profile not found"})
		}
		//Notfound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "Achivement updated!"})

}

/*
	プライバシー関連
*/

// プライバシー情報を取得
func GetPrivacyProfile(c echo.Context) error {
	//useridをヘッダーから取得
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)

	//取得
	privacy, err := services.GetPrivacyProfileService(userID)

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
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)
	var req services.PrivacyInput

	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdatePrivacyProfileById(userID, req); err != nil {
		//Notfoundのとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "privacy profile not found"})
		}
		//Notfound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}

	return c.JSON(http.StatusOK, echo.Map{"status": "privacy updated!"})
}

/*
	地域情報関連
*/

// 地域情報の取得
func GetRegionProfile(c echo.Context) error {

	//useridをヘッダーから取得
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)

	//取得
	region, err := services.GetRegionProfileService(userID)

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
	// userID := c.Request().Header.Get("userid")
	userID := c.Get("UserID").(string)

	//格納用の地域情報
	var req struct {
		RegionID string `json:"region_id"`
	}
	//リクエスト整形(bind)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	//エラー処理
	if err := services.UpdateRegionById(userID, req.RegionID); err != nil {
		//NotFOundのとき
		if errors.Is(err,gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "region profile not found"})
		}
		//NotFound以外のエラー
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "update failed"})
	}
	
	return c.JSON(http.StatusOK, echo.Map{"status": "region profile updated!"})
}
