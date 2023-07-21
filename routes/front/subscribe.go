package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SubscribeRoutes(front *gin.RouterGroup) {

	subscribeRoutes := front.Group("/subscribes")
	{
		// AddLike -> customer bir pro registrasiya etmek ucin ulanylyar
		subscribeRoutes.POST("", middlewares.CheckCustomer(), controllers.AddOrRemoveSubscribe)

		// // GetCustomerLikes funksiya frontdan token bar bolan yagdayynda
		// // musderinin halanlarym sahypasyna gosan harytlaryny getiryar
		// likeRoutes.GET("", middlewares.CheckCustomer(), controllers.GetCustomerLikes)

	}

}
