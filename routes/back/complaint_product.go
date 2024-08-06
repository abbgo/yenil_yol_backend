package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ComplaintProductRoutes(back *gin.RouterGroup) {
	complaintProductRoutes := back.Group("/complaint-products").Use(middlewares.CheckTokenAdminOrShopOwner())
	{
		// GetComplaintProduct - funksiya sikayatly harytlary admin tarapda gorkezmek ucin
		complaintProductRoutes.GET("", controllers.GetComplaintProducts)

		// GetProductComplaints - funksiya haryda edilen sikayatlary admin tarapda gorkezmek ucin
		complaintProductRoutes.GET(":product_id", controllers.GetProductComplaints)
	}
}
