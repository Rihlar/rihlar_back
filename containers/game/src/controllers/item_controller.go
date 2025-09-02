package controllers

import (
	"game/logger"
	"game/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

var itemService = services.ItemService{} // サービスの実体を作る。

// 所持アイテム取得
func GetItemBoxHandler(ctx echo.Context) error {
	// ユーザーの特定する TODO:
	// id := ctx.Get("UserID").(string)
	id := ctx.Request().Header.Get("UserID")

	// サービスに渡す
	itemBox, err := itemService.GetItemBox(id)
	if err != nil {
		logger.PrintErr("アイテム取得エラー", itemBox)
	}

		// 成功ログ
	logger.Println("Successful itembox get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": itemBox,
	})
}

func GetItemDeteileHandler(ctx echo.Context) error {
	// itemID取得
	itemId := ctx.Request().Header.Get("ItemID")

	// サービスに渡す
	item, err := itemService.GetItemDeteile(itemId)
	if err != nil {
		logger.PrintErr("アイテム取得エラー", err)
		return err
	}

	// 成功ログ
	logger.Println("Successful item get.")

	// レスポンス
	return ctx.JSON(http.StatusOK, echo.Map{
		"Data": item,
	})
}