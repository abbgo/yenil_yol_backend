package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SettingRoutes(back *gin.RouterGroup) {
	backSettingApi := back.Group("/setting").Use(middlewares.CheckToken("admin"))
	{
		// CreateSetting -> Setting gosmak ulanylar
		backSettingApi.POST("", controllers.CreateSetting)

		// UpdateSetting -> setting - in maglumatlaryny update etmek ucin ulanylyar
		backSettingApi.PUT("", controllers.UpdateSetting)
	}
}
