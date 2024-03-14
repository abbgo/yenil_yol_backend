package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func DimensionGroupRoutes(back *gin.RouterGroup) {
	backDimensionGroupApi := back.Group("/dimension-groups").Use(middlewares.CheckToken("admin"))
	{
		// CreateDimensionGroup -> DimensionGroup gosmak ucin ulanylar
		backDimensionGroupApi.POST("", controllers.CreateDimensionGroup)

		// UpdateDimensionGroup -> id boyunca DimensionGroup - in maglumatlaryny update etmek ucin ulanylyar
		backDimensionGroupApi.PUT("", controllers.UpdateDimensionGroup)

		// GetDimensionGroupByID -> id - si boyunca DimensionGroup - yn maglumatlaryny almak ucin ulanylyar
		backDimensionGroupApi.GET(":id", controllers.GetDimensionGroupByID)

		// GetDimensionGroups -> Ahli DimensionGroup - laryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backDimensionGroupApi.GET("", controllers.GetDimensionGroups)

		// DeleteDimensionGroupByID -> id boyunca DimensionGroup - y korzina salmak ucin ulanylyar
		backDimensionGroupApi.DELETE(":id", controllers.DeleteDimensionGroupByID)

		// RestoreDimensionGroupByID -> id boyunca DimensionGroup - y korzinadan cykarmak ucin ulanylyar
		backDimensionGroupApi.GET(":id/restore", controllers.RestoreDimensionGroupByID)

		// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
		// backDimensionGroupApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)
	}
}
