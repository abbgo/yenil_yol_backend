package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func TranslationPagesRoutes(back *gin.RouterGroup) {

	backPageTrApi := back.Group("/page-translations")
	{
		{

			// CreateTranslationPage -> TranslationPage gosmak ulanylar
			backPageTrApi.POST("", controllers.CreatePageTr)

			// // UpdatePageByID -> id boyunca Page - in maglumatlaryny update etmek ucin ulanylyar
			// backPageTrApi.PUT("", controllers.UpdatePageByID)

			// // GetPageByID -> id - si boyunca pAGE - in maglumatlaryny almak ucin ulanylyar
			// backPageTrApi.GET(":id", controllers.GetPageByID)

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
