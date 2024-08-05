package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ComplaintProductRoutes(back *gin.RouterGroup) {
	complaintProductRoutes := back.Group("/complaint-products").Use(middlewares.CheckTokenAdminOrShopOwner())
	{
		// GetComplaintProduct - funksiya haryda edilen sikayatlary admin tarapda gorkezmek ucin
		complaintProductRoutes.GET("", controllers.GetComplaintProduct)
	}
}
