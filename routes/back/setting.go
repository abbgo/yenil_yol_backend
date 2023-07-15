package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func SettingRoutes(back *gin.RouterGroup) {

	backSettingApi := back.Group("/setting")
	{
		{

			// CreateSetting -> Setting gosmak ulanylar
			backSettingApi.POST("", controllers.CreateSetting)

			// UpdateSetting -> setting - in maglumatlaryny update etmek ucin ulanylyar
			backSettingApi.PUT("", controllers.UpdateSetting)

			// // GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
			// backSettingApi.GET(":id", controllers.GetCategoryByID)

			// // GetCategories -> Ahli Category - leryn maglumatlaryny request query - den gelen
			// // limit we page boyunca pagination ulanyp almak ucin ulanylyar
			// backSettingApi.GET("", controllers.GetCategories)

			// // DeleteCategoryByID -> id boyunca category - ni korzina salmak ucin ulanylyar
			// backSettingApi.DELETE(":id", controllers.DeleteCategoryByID)

			// // RestoreCategoryByID -> id boyunca category - ni korzinadan cykarmak ucin ulanylyar
			// backSettingApi.GET(":id/restore", controllers.RestoreCategoryByID)

			// // DeletePermanentlyCategoryByID -> id boyunca category - i doly (korzinadan) pozmak ucin ulanylyar
			// backSettingApi.DELETE(":id/delete", controllers.DeletePermanentlyCategoryByID)

		}
	}

}
