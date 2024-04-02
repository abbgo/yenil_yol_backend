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
	}

}
