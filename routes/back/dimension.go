package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func DimensionRoutes(back *gin.RouterGroup) {
	backDimensionApi := back.Group("/dimensions").Use(middlewares.CheckToken("admin"))
	{
		// CreateDimensionGroup -> DimensionGroup gosmak ucin ulanylar
		backDimensionApi.POST("", controllers.CreateDimension)

		// UpdateDimensionGroup -> id boyunca DimensionGroup - in maglumatlaryny update etmek ucin ulanylyar
		// backDimensionApi.PUT("", controllers.UpdateDimensionGroup)

		// // GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
		// backDimensionApi.GET(":id", controllers.GetBrendByID)

		// // GetBrends -> Ahli Brend - leryn maglumatlaryny request query - den gelen
		// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
		// backDimensionApi.GET("", controllers.GetBrends)

		// // DeleteBrendByID -> id boyunca brend - i korzina salmak ucin ulanylyar
		// backDimensionApi.DELETE(":id", controllers.DeleteBrendByID)

		// // RestoreBrendByID -> id boyunca brend - i korzinadan cykarmak ucin ulanylyar
		// backDimensionApi.GET(":id/restore", controllers.RestoreBrendByID)

		// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
		// backDimensionApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)
	}
}
