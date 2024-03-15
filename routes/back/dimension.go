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

		// UpdateDimension -> id boyunca Dimension - yn maglumatlaryny update etmek ucin ulanylyar
		backDimensionApi.PUT("", controllers.UpdateDimension)

		// GetDimensionByID -> id - si boyunca Dimension - yn maglumatlaryny almak ucin ulanylyar
		backDimensionApi.GET(":id", controllers.GetDimensionByID)

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
