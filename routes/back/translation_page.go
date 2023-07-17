package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func TranslationPagesRoutes(back *gin.RouterGroup) {

	backTrpageApi := back.Group("/translation-page")
	{
		{

			// CreatePage -> Page gosmak ulanylar
			backTrpageApi.POST("", controllers.CreatePage)

			// // UpdatePageByID -> id boyunca Page - in maglumatlaryny update etmek ucin ulanylyar
			// backTrpageApi.PUT("", controllers.UpdatePageByID)

			// // GetPageByID -> id - si boyunca pAGE - in maglumatlaryny almak ucin ulanylyar
			// backTrpageApi.GET(":id", controllers.GetPageByID)

			// // GetPages -> Ahli Page - leryn maglumatlaryny request query - den gelen
			// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// backTrpageApi.GET("", controllers.GetPages)

			// // DeletePageByID -> id boyunca page - i korzina salmak ucin ulanylyar
			// backTrpageApi.DELETE(":id", controllers.DeletePageByID)

			// // RestorePageByID -> id boyunca page - i korzinadan cykarmak ucin ulanylyar
			// backTrpageApi.GET(":id/restore", controllers.RestorePageByID)

			// // DeletePermanentlyPageByID -> id boyunca page - i doly (korzinadan) pozmak ucin ulanylyar
			// backTrpageApi.DELETE(":id/delete", controllers.DeletePermanentlyPageByID)

		}
	}

}
