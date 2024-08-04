package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ComplaintRoutes(front *gin.RouterGroup) {
	complaintRoutes := front.Group("/complaints")
	{
		// GetComplaints funksiya fronta sikayatlary ugratmak ucin ulanylyar
		complaintRoutes.GET("", controllers.GetComplaints)
	}
}
