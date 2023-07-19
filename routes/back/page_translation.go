package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func TranslationPagesRoutes(back *gin.RouterGroup) {

	backPageTrApi := back.Group("/page-translations")
	{
		{

			// CreateTranslationPage -> PageTranslation gosmak ulanylar
			backPageTrApi.POST("", controllers.CreatePageTr)

			// UpdatePageTrByID -> id boyunca PageTranslation - in maglumatlaryny update etmek ucin ulanylyar
			backPageTrApi.PUT("", controllers.UpdatePageTrByID)

			// GetPageTrByID -> id - si boyunca PageTranslation - in maglumatlaryny almak ucin ulanylyar
			backPageTrApi.GET(":id", controllers.GetPageTrByID)

			// // GetPages -> Ahli Page - leryn maglumatlaryny request query - den gelen
			// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// backPageTrApi.GET("", controllers.GetPages)

			// // DeletePageByID -> id boyunca page - i korzina salmak ucin ulanylyar
			// backPageTrApi.DELETE(":id", controllers.DeletePageByID)

			// // RestorePageByID -> id boyunca page - i korzinadan cykarmak ucin ulanylyar
			// backPageTrApi.GET(":id/restore", controllers.RestorePageByID)

			// // DeletePermanentlyPageByID -> id boyunca page - i doly (korzinadan) pozmak ucin ulanylyar
			// backPageTrApi.DELETE(":id/delete", controllers.DeletePermanentlyPageByID)

		}
	}

}
