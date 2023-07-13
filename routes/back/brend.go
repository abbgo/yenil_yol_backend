package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func BrendRoutes(back *gin.RouterGroup) {

	backBrendApi := back.Group("/brends")
	{
		{

			// surat gosmak ya-da uytgetmek ucin ulanylyar
			backBrendApi.POST("", controllers.AddOrUpdateImage)

			// surat pozmak ucin ulanylyar
			backBrendApi.DELETE("", controllers.DeleteImage)

		}
	}

}
