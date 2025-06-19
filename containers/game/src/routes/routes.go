package routes

// エンドポイントのルーティング
import (
	"game/controllers"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ルーティング　nagasugi kimosugi yurusenai
func InitRoutes(e *echo.Echo, db *gorm.DB) {

	//
	e.GET("/users/:id", controllers.GetMyRankingHandler)

}
