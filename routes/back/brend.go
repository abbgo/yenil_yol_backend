package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func BrendRoutes(back *gin.RouterGroup) {

	backBrendApi := back.Group("/brends")
	{
		{

			// CreateBrend -> Brend gosmak ulanylar
			backBrendApi.POST("", controllers.CreateBrend)

			// UpdateBrendByID -> id boyunca Brend - in maglumatlaryny update etmek ucin ulanylyar
			backBrendApi.PUT("", controllers.UpdateBrendByID)

			// GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
			backBrendApi.GET(":id", controllers.GetBrendByID)

			// // GetShops -> Ahli Shop - laryn maglumatlaryny request query - den gelen
			// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// // eger request query - den shop_owner_id gelse sol shop_owner degisli
			// // shop - laryn maglumatlary alynyar
			// backBrendApi.GET("", controllers.GetShops)

			// // DeleteShopByID -> id boyunca shop - y korzina salmak ucin ulanylyar
			// backBrendApi.DELETE(":id", controllers.DeleteShopByID)

			// // RestoreShopByID -> id boyunca shop - y korzinadan cykarmak ucin ulanylyar
			// backBrendApi.GET(":id/restore", controllers.RestoreShopByID)

			// // DeletePermanentlyShopByID -> id boyunca shop - y doly (korzinadan) pozmak ucin ulanylyar
			// backBrendApi.DELETE(":id/delete", controllers.DeletePermanentlyShopByID)

		}
	}

}
