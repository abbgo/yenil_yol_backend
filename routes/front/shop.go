package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ShopRoutes(front *gin.RouterGroup) {
	shopRoutes := front.Group("/shops")
	{
		// // AddOrRemoveSubscribe -> customer -e like gosmak ya-da pozmak ucin ulanylyar
		// // gosmak ucin request query - de status = true
		// // pozmak ucin request query - de status = false ugratmaly
		// // bu api - yn islemegi ucin customer token gerek
		// likeRoutes.POST("", middlewares.CheckCustomer(), controllers.AddOrRemoveLike)

		// // GetCustomerLikes funksiya frontdan token bar bolan yagdayynda
		// // musderinin halanlarym sahypasyna gosan harytlaryny getiryar
		// likeRoutes.GET("", middlewares.CheckCustomer(), controllers.GetCustomerLikes)

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
	}
}
