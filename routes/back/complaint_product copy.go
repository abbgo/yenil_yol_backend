package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ComplaintProductRoutes(back *gin.RouterGroup) {
	complaintProductRoutes := back.Group("/complaint-products")
	{
		// GetComplaintProduct - funksiya haryda edilen sikayatlary admin tarapda gorkezmek ucin
		complaintProductRoutes.GET("", controllers.GetComplaintProduct)
	}
}
