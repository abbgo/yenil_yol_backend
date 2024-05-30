package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ShopRoutes(front *gin.RouterGroup) {
	shopRoutes := front.Group("/shops")
	{
		// GetShops -> Ahli Shop - laryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		// eger request query - den shop_owner_id gelse sol shop_owner degisli
		// shop - laryn maglumatlary alynyar
		shopRoutes.GET("", controllers.GetShops)

		// GetShopByID -> Dine bir Shop - yn maglumatlaryny request param - dan gelen
		// id boyunca alynyar
		shopRoutes.GET(":id", controllers.GetShopByID)

		// GetShopsForMap - funksiya front - daky karta ucin gerek bolan ahli magazynlary alyar
		shopRoutes.GET("map", controllers.GetShopsForMap)

		// GetShopByIDs - funksiya query - dan gelen id - ler boyunca magazynlary alyar
		shopRoutes.GET(":id", controllers.GetShopByIDs)
	}
}
