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
			backApi.GET("", controllers.GetShops)
			backApi.DELETE(":id", controllers.DeleteShopByID)
			backApi.GET(":id/restore", controllers.RestoreShopByID)
			backApi.DELETE(":id/delete", controllers.DeletePermanentlyShopByID)

		}
	}

	back.POST("image", controllers.AddOrUpdateImage)
	back.DELETE("image", controllers.DeleteImage)

}
