package controllers

import (
	"gcore/services"

	"github.com/labstack/echo/v4"
)

func ReportMovement(ctx echo.Context) error {
	// body を取得する
	req := new(services.MovementReportRequest)
	if err := ctx.Bind(req); err != nil {
		return err
	}

	// サービスを呼び出す
	

	return nil
}