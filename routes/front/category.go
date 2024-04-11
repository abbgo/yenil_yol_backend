package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(front *gin.RouterGroup) {
	categoryRoutes := front.Group("/categories")
	{
		// GetCategoriesShopID - request parameter-den gelen shop_id boyunca category - leri ugratyar
		categoryRoutes.GET(":shop_id", controllers.GetCategoriesShopID)

		// GetCategoriesByCategoryID - request parameter-den gelen shop_id we category_id boyunca category - leri ugratyar
		categoryRoutes.GET(":shop_id/:category_id", controllers.GetCategoriesByCategoryID)
	}
}
