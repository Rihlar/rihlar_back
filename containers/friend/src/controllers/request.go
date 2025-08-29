package controllers

import (
	"friend/logger"
	"friend/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 受信済みリクエストを取得する関数
func GetRecvedRequest(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	data, err := services.GetRecvedRequest(userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get recved request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}

// リクエストを拒否する関数
func RejectRequest(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// 拒否対象のユーザーID を取得
	targetUserId := ctx.Request().Header.Get("userId")

	// サービスに渡す
	err := services.RejectRequest(userId, targetUserId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to reject request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"result": "success"})
}

// リクエストを承認する関数
func AcceptRequest(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// 承認対象のユーザーID を取得
	targetUserId := ctx.Request().Header.Get("userId")

	// サービスに渡す
	err := services.AcceptRequest(userId, targetUserId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to accept request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"result": "success"})
}

// フレンドリクエストをキャンセルする関数
func CancelRequest(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// キャンセル対象のユーザーID を取得
	targetUserId := ctx.Request().Header.Get("userId")

	// サービスに渡す
	err := services.CancelRequest(userId, targetUserId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to cancel request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"result": "success"})
}

// 送信済みリクエスト取得する関数
func GetSentRequest(ctx echo.Context) error {
	// ユーザーID を取得
	userId := ctx.Get("UserID").(string)

	// サービスに渡す
	data, err := services.GetSentRequest(userId)

	// エラー処理
	if err != nil {
		logger.PrintErr(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get sent request","message" : err.Error()})
	}

	return ctx.JSON(http.StatusOK, data)
}