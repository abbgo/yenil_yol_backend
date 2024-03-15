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

		// GetDimensionsByGroupID -> dimension_group_id boyunca Dimension - laryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backDimensionApi.GET("", controllers.GetDimensionsByGroupID)

		// DeleteDimensionByID -> id boyunca dimension - y korzina salmak ucin ulanylyar
		backDimensionApi.DELETE(":id", controllers.DeleteDimensionByID)

		// RestoreDimensionByID -> id boyunca dimension - y korzinadan cykarmak ucin ulanylyar
		backDimensionApi.GET(":id/restore", controllers.RestoreDimensionByID)

		// DeletePermanentlyDimensionByID -> id boyunca dimension - y doly (korzinadan) pozmak ucin ulanylyar
		backDimensionApi.DELETE(":id/delete", controllers.DeletePermanentlyDimensionByID)
	}
}
