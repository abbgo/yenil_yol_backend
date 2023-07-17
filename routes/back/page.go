package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func PagesRoutes(back *gin.RouterGroup) {

	backPageApi := back.Group("/pages")
	{
		{

			// CreatePage -> Page gosmak ulanylar
			backPageApi.POST("", controllers.CreatePage)

			// UpdatePageByID -> id boyunca Page - in maglumatlaryny update etmek ucin ulanylyar
			backPageApi.PUT("", controllers.UpdatePageByID)

			// GetPageByID -> id - si boyunca pAGE - in maglumatlaryny almak ucin ulanylyar
			backPageApi.GET(":id", controllers.GetPageByID)

			// GetPages -> Ahli Page - leryn maglumatlaryny request query - den gelen
			// limit we page boyunca pagination ulanyp almak ucin ulanylyar
			backPageApi.GET("", controllers.GetPages)

			// DeletePageByID -> id boyunca page - i korzina salmak ucin ulanylyar
			backPageApi.DELETE(":id", controllers.DeletePageByID)

			// RestorePageByID -> id boyunca page - i korzinadan cykarmak ucin ulanylyar
			backPageApi.GET(":id/restore", controllers.RestorePageByID)

			// DeletePermanentlyPageByID -> id boyunca page - i doly (korzinadan) pozmak ucin ulanylyar
			backPageApi.DELETE(":id/delete", controllers.DeletePermanentlyPageByID)

		}
	}

}
