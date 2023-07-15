package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func BrendCategories(back *gin.RouterGroup) {

	backCategoryApi := back.Group("/categories")
	{
		{

			// CreateCategory -> Category gosmak ulanylar
			backCategoryApi.POST("", controllers.CreateCategory)

			// UpdateCategoryByID -> id boyunca Category - in maglumatlaryny update etmek ucin ulanylyar
			backCategoryApi.PUT("", controllers.UpdateCategoryByID)

			// GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
			backCategoryApi.GET(":id", controllers.GetCategoryByID)

			// // GetBrends -> Ahli Brend - leryn maglumatlaryny request query - den gelen
			// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// backCategoryApi.GET("", controllers.GetBrends)

			// // DeleteBrendByID -> id boyunca brend - i korzina salmak ucin ulanylyar
			// backCategoryApi.DELETE(":id", controllers.DeleteBrendByID)

			// // RestoreBrendByID -> id boyunca brend - i korzinadan cykarmak ucin ulanylyar
			// backCategoryApi.GET(":id/restore", controllers.RestoreBrendByID)

			// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
			// backCategoryApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)

		}
	}

}
