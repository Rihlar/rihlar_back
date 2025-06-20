package routes

// エンドポイントのルーティング
import (
	"game/controllers"
	"github.com/labstack/echo/v4"
)

// ルーティング　
func InitRoutes(e *echo.Echo) *echo.Echo {

	//ranking取得
	e.GET("/ranking/personal/:user_id", controllers.GetMyRankingHandler)

	return e
}
