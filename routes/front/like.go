package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(front *gin.RouterGroup) {

	likeRoutes := front.Group("/likes")
	{
		// AddLike -> customer bir pro registrasiya etmek ucin ulanylyar
		likeRoutes.POST("", middlewares.CheckCustomer(), controllers.AddOrRemoveLike)

	}

}
