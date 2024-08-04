package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ComplaintProductRoutes(front *gin.RouterGroup) {
	complaintProductRoutes := front.Group("/complaint-products")
	{
		// GetCustomerLikes funksiya frontdan token bar bolan yagdayynda
		// musderinin halanlarym sahypasyna gosan harytlaryny getiryar
		complaintProductRoutes.GET("", controllers.GetComplaints)
	}
}
