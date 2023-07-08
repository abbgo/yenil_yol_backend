package shopOwnerApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/admin"

	"github.com/gin-gonic/gin"
)

func ShopOwnerRoutes(back *gin.RouterGroup) {

	shopOwner := back.Group("/shop-owner")
	{
		{
			// RegisterAdmin admin registrasiya etmek ucin ulanylyar.
			// Admini dine super admin registrasiya edip bilyar. Admin admin registrasiya edip bilenok
			shopOwner.POST("register", controllers.RegisterShopOwner)
			// admin.POST("register", middlewares.IsSuperAdmin(), adminController.RegisterAdmin)

		}
	}

}
