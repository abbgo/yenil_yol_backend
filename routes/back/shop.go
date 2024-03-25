package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func BackShopOwnerRoutes(back *gin.RouterGroup) {
	backShopOwnerApi := back.Group("/shops").Use(middlewares.CheckTokenAdminOrCustomer())
	{
		// CreateShop -> gosmak ulanylar
		backShopOwnerApi.POST("", controllers.CreateShop)

		// UpdateShopByID -> id boyunca Shop - yn maglumatlaryny update etmek ucin ulanylyar
		backShopOwnerApi.PUT("", controllers.UpdateShopByID)

		// GetShopByID -> id - si boyunca Shop - yn maglumatlaryny almak ucin ulanylyar
		backShopOwnerApi.GET(":id", controllers.GetShopByID)

		// GetShops -> Ahli Shop - laryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		// eger request query - den shop_owner_id gelse sol shop_owner degisli
		// shop - laryn maglumatlary alynyar
		backShopOwnerApi.GET("", controllers.GetShops)

		// DeleteShopByID -> id boyunca shop - y korzina salmak ucin ulanylyar
		backShopOwnerApi.DELETE(":id", controllers.DeleteShopByID)

		// RestoreShopByID -> id boyunca shop - y korzinadan cykarmak ucin ulanylyar
		backShopOwnerApi.GET(":id/restore", controllers.RestoreShopByID)

		// DeletePermanentlyShopByID -> id boyunca shop - y doly (korzinadan) pozmak ucin ulanylyar
		backShopOwnerApi.DELETE(":id/delete", controllers.DeletePermanentlyShopByID)
	}
}
