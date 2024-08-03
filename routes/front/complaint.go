package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func ComplaintRoutes(front *gin.RouterGroup) {

	complaintRoutes := front.Group("/complaints")
	{
		// // AddOrRemoveLike -> customer -e like gosmak ya-da pozmak ucin ulanylyar
		// // gosmak ucin request query - de status = true
		// // pozmak ucin request query - de status = false ugratmaly
		// // bu api - yn islemegi ucin customer token gerek
		// complaintRoutes.POST("", middlewares.CheckToken("customer"), controllers.AddOrRemoveLike)

		// GetCustomerLikes funksiya frontdan token bar bolan yagdayynda
		// musderinin halanlarym sahypasyna gosan harytlaryny getiryar
		complaintRoutes.GET("", controllers.GetComplaints)

	}

}
