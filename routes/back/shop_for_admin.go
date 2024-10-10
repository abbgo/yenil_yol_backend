package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func BackShopForAdminRoutes(back *gin.RouterGroup) {
	backShopOwnerApi := back.Group("/shops/admin").Use(middlewares.CheckToken("admin"))
	{
		// GetAdminShops -> Ahli Shop - laryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backShopOwnerApi.GET("", controllers.GetAdminShops)

		// UpdateShopCreatedStatus -> Funksiya shop - yn created status - yny uytgetmek
		// ucin ulanylyar
		backShopOwnerApi.PUT("created-status", controllers.UpdateShopCreatedStatus)

		// UpdateShopBrandStatus -> Funksiya shop - yn ofisialny dukan yagdayyny uytgetmek
		// ucin ulanylyar
		backShopOwnerApi.PUT("brand-status", controllers.UpdateShopBrandStatus)
	}
}
