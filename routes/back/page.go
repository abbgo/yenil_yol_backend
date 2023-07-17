package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func PagesRoutes(back *gin.RouterGroup) {

	backBrendApi := back.Group("/pages")
	{
		{

			// CreatePage -> Page gosmak ulanylar
			backBrendApi.POST("", controllers.CreatePage)

			// UpdatePageByID -> id boyunca Page - in maglumatlaryny update etmek ucin ulanylyar
			backBrendApi.PUT("", controllers.UpdatePageByID)

			// // GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
			// backBrendApi.GET(":id", controllers.GetBrendByID)

			// // GetBrends -> Ahli Brend - leryn maglumatlaryny request query - den gelen
			// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// backBrendApi.GET("", controllers.GetBrends)

			// // DeleteBrendByID -> id boyunca brend - i korzina salmak ucin ulanylyar
			// backBrendApi.DELETE(":id", controllers.DeleteBrendByID)

			// // RestoreBrendByID -> id boyunca brend - i korzinadan cykarmak ucin ulanylyar
			// backBrendApi.GET(":id/restore", controllers.RestoreBrendByID)

			// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
			// backBrendApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)

		}
	}

}
