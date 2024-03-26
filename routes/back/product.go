package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(back *gin.RouterGroup) {
	backProductApi := back.Group("/products").Use(middlewares.CheckTokenAdminOrCustomer())
	{
		// // CreateProduct -> Product gosmak ulanylar
		backProductApi.POST("", controllers.CreateProduct)

		// UpdateProductByID -> id boyunca Product - in maglumatlaryny update etmek ucin ulanylyar
		backProductApi.PUT("", controllers.UpdateProductByID)

		// GetProductByID -> id - si boyunca Product - yn maglumatlaryny almak ucin ulanylyar
		backProductApi.GET(":id", controllers.GetProductByID)

		// GetProducts -> Ahli Product - leryn maglumatlaryny request query - den gelen
		// limit we page boyunca pagination ulanyp almak ucin ulanylyar
		backProductApi.GET("", controllers.GetProducts)

		// DeleteProductByID -> id boyunca product - y korzina salmak ucin ulanylyar
		backProductApi.DELETE(":id", controllers.DeleteProductByID)

		// RestoreProductByID -> id boyunca product - y korzinadan cykarmak ucin ulanylyar
		backProductApi.GET(":id/restore", controllers.RestoreProductByID)

		// DeletePermanentlyProductByID -> id boyunca product - y doly (korzinadan) pozmak ucin ulanylyar
		backProductApi.DELETE(":id/delete", controllers.DeletePermanentlyProductByID)
	}
}
