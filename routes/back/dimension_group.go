package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func DimensionGroupRoutes(back *gin.RouterGroup) {
	backDimensionGroupApi := back.Group("/dimension-groups").Use(middlewares.CheckTokenAdminOrShopOwner())
	{
		// CreateDimensionGroup -> DimensionGroup gosmak ucin ulanylar
		backDimensionGroupApi.POST("", controllers.CreateDimensionGroup)

		// UpdateDimensionGroup -> id boyunca DimensionGroup - in maglumatlaryny update etmek ucin ulanylyar
		backDimensionGroupApi.PUT("", controllers.UpdateDimensionGroup)

		// GetDimensionGroupByID -> id - si boyunca DimensionGroup - yn maglumatlaryny almak ucin ulanylyar
		backDimensionGroupApi.GET(":id", controllers.GetDimensionGroupByID)

		// GetDimensionGroupsWithDimensions -> Ahli DimensionGroup - laryn maglumatlaryny dimension -lary bilen
		// request query - den gelen limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backDimensionGroupApi.GET("", controllers.GetDimensionGroupsWithDimensions)

		// GetDimensionGroupsWithDimensionsList -> Ahli DimensionGroup - laryn maglumatlaryny dimension - lary ( array edip ) bilen
		// request query - den gelen limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backDimensionGroupApi.GET("with-dimensions", controllers.GetDimensionGroupsWithDimensionsList)

		// DeleteDimensionGroupByID -> id boyunca DimensionGroup - y korzina salmak ucin ulanylyar
		backDimensionGroupApi.DELETE(":id", controllers.DeleteDimensionGroupByID)

		// RestoreDimensionGroupByID -> id boyunca DimensionGroup - y korzinadan cykarmak ucin ulanylyar
		backDimensionGroupApi.GET(":id/restore", controllers.RestoreDimensionGroupByID)

		// DeletePermanentlyDimensionGroupByID -> id boyunca DimensionGroup - y doly (korzinadan) pozmak ucin ulanylyar
		backDimensionGroupApi.DELETE(":id/delete", controllers.DeletePermanentlyDimensionGroupByID)
	}
}
