package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(back *gin.RouterGroup) {

	backProductApi := back.Group("/products")
	{
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

			// // RestoreBrendByID -> id boyunca brend - i korzinadan cykarmak ucin ulanylyar
			// backProductApi.GET(":id/restore", controllers.RestoreBrendByID)

			// // DeletePermanentlyBrendByID -> id boyunca brend - i doly (korzinadan) pozmak ucin ulanylyar
			// backProductApi.DELETE(":id/delete", controllers.DeletePermanentlyBrendByID)

		}
	}

}
