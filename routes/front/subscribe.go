package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SubscribeRoutes(front *gin.RouterGroup) {

	subscribeRoutes := front.Group("/subscribes")
	{
		// AddOrRemoveSubscribe -> customer -e subscribe shop gosmak ya-da pozmak ucin ulanylyar
		// gosmak ucin request query - de status = true
		// pozmak ucin request query - de status = false ugratmaly
		// by api - nin islemegi ucin customer token gerek
		subscribeRoutes.POST("", middlewares.CheckCustomer(), controllers.AddOrRemoveSubscribe)

		// GetCustomerSubscribes funksiya frontdan token bar bolan yagdayynda
		// musderinin subscribe sahypasyna gosan shop - laryny getiryar
		subscribeRoutes.GET("", middlewares.CheckCustomer(), controllers.GetCustomerSubscribes)

	}

}
