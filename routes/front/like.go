package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(front *gin.RouterGroup) {

	likeRoutes := front.Group("/likes")
	{
		// AddOrRemoveSubscribe -> customer -e like gosmak ya-da pozmak ucin ulanylyar
		// gosmak ucin request query - de status = true
		// pozmak ucin request query - de status = false ugratmaly
		// bu api - yn islemegi ucin customer token gerek
		likeRoutes.POST("", middlewares.CheckToken("customer"), controllers.AddOrRemoveLike)

		// GetCustomerLikes funksiya frontdan token bar bolan yagdayynda
		// musderinin halanlarym sahypasyna gosan harytlaryny getiryar
		likeRoutes.GET("", middlewares.CheckToken("customer"), controllers.GetCustomerLikes)

	}

}
