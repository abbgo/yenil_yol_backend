package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(front *gin.RouterGroup) {
	productRoutes := front.Group("/products")
	{
		// GetProductByID -> Dine bir Product - yn maglumatlaryny request param - dan gelen
		// id boyunca alynyar
		productRoutes.GET(":id", controllers.GetProductByID)

		// GetProducts -> Product - laryn maglumatlaryny query - den gelen maglumatlar boyunca
		// almak ucin ulanylyar
		productRoutes.GET("", controllers.GetProducts)

		// GetSimilarProductsByCategoryID -> Product - laryn maglumatlaryny query - den gelen maglumatlar boyunca
		// almak ucin ulanylyar
		productRoutes.GET("similars", controllers.GetSimilarProductsByProductID)

		//GetProductsByIDs - funksiya query - dan gelen id - ler boyunca harytlary alyar
		productRoutes.GET("favorite", controllers.GetProductsByIDs)
	}
}
