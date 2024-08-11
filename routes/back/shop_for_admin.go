package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func BackShopForAdminRoutes(back *gin.RouterGroup) {
	backShopOwnerApi := back.Group("/shops/admin").Use(middlewares.CheckTokenAdminOrShopOwner())
	{
		// GetShops -> Ahli Shop - laryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		// eger request query - den shop_owner_id gelse sol shop_owner degisli
		// shop - laryn maglumatlary alynyar
		backShopOwnerApi.GET("", controllers.GetShops)
	}
}
