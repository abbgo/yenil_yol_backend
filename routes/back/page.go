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

			// GetPageByID -> id - si boyunca pAGE - in maglumatlaryny almak ucin ulanylyar
			backBrendApi.GET(":id", controllers.GetPageByID)

			// GetPages -> Ahli Page - leryn maglumatlaryny request query - den gelen
			// limit we page boyunca pagination ulanyp almak ucin ulanylyar
			backBrendApi.GET("", controllers.GetPages)

			// DeletePageByID -> id boyunca page - i korzina salmak ucin ulanylyar
			backBrendApi.DELETE(":id", controllers.DeletePageByID)

			// RestorePageByID -> id boyunca page - i korzinadan cykarmak ucin ulanylyar
			backBrendApi.GET(":id/restore", controllers.RestorePageByID)

			// DeletePermanentlyPageByID -> id boyunca page - i doly (korzinadan) pozmak ucin ulanylyar
			backBrendApi.DELETE(":id/delete", controllers.DeletePermanentlyPageByID)

		}
	}

}
