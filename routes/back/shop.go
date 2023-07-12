package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func BackRoutes(back *gin.RouterGroup) {

	backApi := back.Group("/shops")
	{
		{
			// CreateShop -> gosmak ulanylar
			backApi.POST("", controllers.CreateShop)

			// UpdateShopByID -> id boyunca Shop - yn maglumatlaryny update etmek ucin ulanylyar
			backApi.PUT("", controllers.UpdateShopByID)

			// GetShopByID -> id - si boyunca Shop - yn maglumatlaryny almak ucin ulanylyar
			backApi.GET(":id", controllers.GetShopByID)

			// GetShops -> Ahli Shop - laryn maglumatlaryny request query - den gelen
			// limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// eger request query - den shop_owner_id gelse sol shop_owner degisli
			// shop - laryn maglumatlary alynyar
			backApi.GET("", controllers.GetShops)

			// DeleteShopByID -> id boyunca shop - y korzina salmak ucin ulanylyar
			backApi.DELETE(":id", controllers.DeleteShopByID)

			// RestoreShopByID -> id boyunca shop - y korzinadan cykarmak ucin ulanylyar
			backApi.GET(":id/restore", controllers.RestoreShopByID)

			// DeletePermanentlyShopByID -> id boyunca shop - y doly (korzinadan) pozmak ucin ulanylyar
			backApi.DELETE(":id/delete", controllers.DeletePermanentlyShopByID)

		}
	}

	back.POST("image", controllers.AddOrUpdateImage)
	back.DELETE("image", controllers.DeleteImage)

}
