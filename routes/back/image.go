package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func BackImageRoutes(back *gin.RouterGroup) {
	backImageApi := back.Group("/image")
	{
		// surat gosmak ya-da uytgetmek ucin ulanylyar
		backImageApi.POST("", controllers.AddOrUpdateImage)

		// surat pozmak ucin ulanylyar
		backImageApi.DELETE("", controllers.DeleteImage)
	}
}
