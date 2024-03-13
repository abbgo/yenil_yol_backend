package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func DimensionGroupRoutes(back *gin.RouterGroup) {

	backDimensionGroupApi := back.Group("/dimension-groups")
	{

		// CreateDimensionGroup -> DimensionGroup gosmak ucin ulanylar
		backDimensionGroupApi.POST("", controllers.CreateDimensionGroup)

		// UpdateDimensionGroup -> id boyunca DimensionGroup - in maglumatlaryny update etmek ucin ulanylyar
		backDimensionGroupApi.PUT("", controllers.UpdateDimensionGroup)

		// // GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
		// backDimensionGroupApi.GET(":id", controllers.GetBrendByID)

		// // GetBrends -> Ahli Brend - leryn maglumatlaryny request query - den gelen
		// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
		// backDimensionGroupApi.GET("", controllers.GetBrends)

		// // DeleteBrendByID -> id boyunca brend - i korzina salmak ucin ulanylyar
		// backDimensionGroupApi.DELETE(":id", controllers.DeleteBrendByID)

		// // RestoreBrendByID -> id boyunca brend - i korzinadan cykarmak ucin ulanylyar
		// backDimensionGroupApi.GET(":id/restore", controllers.RestoreBrendByID)

		// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
		// backDimensionGroupApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)

	}

}
