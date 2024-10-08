package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(front *gin.RouterGroup) {
	categoryRoutes := front.Group("/categories")
	{
		// GetCategoriesShopID - request parameter-den gelen shop_id boyunca category - leri ugratyar
		categoryRoutes.GET("", controllers.GetCategories)
	}
}
