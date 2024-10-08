package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func BackProductRoutesForAdmin(back *gin.RouterGroup) {
	backProductApi := back.Group("/products/admin").Use(middlewares.CheckToken("admin"))
	{
		// GetProducts -> Ahli Product - leryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backProductApi.GET("", controllers.GetAdminProducts)

		// UpdateShopCreatedStatus -> Funksiya shop - yn created status - yny uytgetmek
		// ucin ulanylyar
		backProductApi.PUT("created-status", controllers.UpdateProductCreatedStatus)
	}
}
