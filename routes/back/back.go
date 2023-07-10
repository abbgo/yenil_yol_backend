package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func ShopOwnerRoutes(back *gin.RouterGroup) {

	backApi := back.Group("/shops")
	{
		{
			// Shop gosmak ulanylar
			backApi.POST("", controllers.CreateShop)

		}
	}

}
