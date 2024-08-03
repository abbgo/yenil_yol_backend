package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ComplaintRoutes(back *gin.RouterGroup) {
	backComplaintApi := back.Group("/complaints").Use(middlewares.CheckToken("admin"))
	{
		// CreateComplaint -> sikayatyn tekstini gosmak ulanylar
		backComplaintApi.POST("", controllers.CreateComplaint)

		// UpdateComplaintByID -> id boyunca sikayatyn tekstini update etmek ucin ulanylyar
		backComplaintApi.PUT("", controllers.UpdateComplaintByID)

		// // GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
		// backBrendApi.GET(":id", controllers.GetBrendByID)

		// // GetBrends -> Ahli Brend - leryn maglumatlaryny request query - den gelen
		// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
		// backBrendApi.GET("", controllers.GetBrends)

		// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
		// backBrendApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)
	}
}
