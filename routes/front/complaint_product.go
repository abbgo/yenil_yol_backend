package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ComplaintProductRoutes(front *gin.RouterGroup) {
	complaintProductRoutes := front.Group("/complaint-products")
	{
		// CreateComplaintProduct - funksiya front tarapda klient haryda sikayat doretmek ucin ulanylyar
		complaintProductRoutes.POST("", controllers.CreateComplaintProduct)
	}
}
