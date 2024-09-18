package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func BackCategoryRoutesForAdmin(back *gin.RouterGroup) {
	backCategoryApi := back.Group("/categories").Use(middlewares.CheckToken("admin"))
	{
		// CreateCategory -> Category gosmak ulanylar
		backCategoryApi.POST("", controllers.CreateCategory)

		// UpdateCategoryByID -> id boyunca Category - in maglumatlaryny update etmek ucin ulanylyar
		backCategoryApi.PUT("", controllers.UpdateCategoryByID)

		// GetBrendByID -> id - si boyunca Brend - in maglumatlaryny almak ucin ulanylyar
		backCategoryApi.GET(":id", controllers.GetCategoryByID)

		// GetCategoriesWithChild -> Ahli parent Category - leryn cagalary bilen bile maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backCategoryApi.GET("with-child", controllers.GetCategoriesWithChild)

		// GetDeletedCategories -> Ahli pozulan kategoriyalar alynyar
		backCategoryApi.GET("deleted", controllers.GetDeletedCategories)

		// GetCategories -> Ahli Category - leryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backCategoryApi.GET("", controllers.GetCategories)

		// CheckForDelete - Kategoriya degisli child kategoriya barmy we kategoriya degisli haryt barmy sony barlayar
		// eger bar bolsa onda kategoriyany pozup bolmayar
		backCategoryApi.GET(":id/check-for-delete", controllers.CheckForDelete)

		// GetParentCategory - id boyunca parent category alyar
		backCategoryApi.GET(":id/parent", controllers.GetParentCategory)

		// DeleteCategoryByID -> id boyunca category - ni korzina salmak ucin ulanylyar
		backCategoryApi.DELETE(":id", controllers.DeleteCategoryByID)

		// RestoreCategoryByID -> id boyunca category - ni korzinadan cykarmak ucin ulanylyar
		backCategoryApi.GET(":id/restore", controllers.RestoreCategoryByID)

		// DeletePermanentlyCategoryByID -> id boyunca category - i doly (korzinadan) pozmak ucin ulanylyar
		backCategoryApi.DELETE(":id/delete", controllers.DeletePermanentlyCategoryByID)
	}
}
