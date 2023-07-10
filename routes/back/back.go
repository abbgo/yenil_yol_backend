package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func BackRoutes(back *gin.RouterGroup) {

	backApi := back.Group("/shops")
	{
		{
			// Shop gosmak ulanylar
			backApi.POST("", controllers.CreateShop)
			backApi.PUT("", controllers.UpdateShopByID)
			backApi.GET(":id", controllers.GetShopByID)

		}
	}

	back.POST("image", controllers.AddOrUpdateImage)
	back.DELETE("image", controllers.DeleteImage)

}
