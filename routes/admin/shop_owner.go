package shopOwnerApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/admin"

	"github.com/gin-gonic/gin"
)

func ShopOwnerRoutes(back *gin.RouterGroup) {

	shopOwner := back.Group("/shop-owner")
	{
		{
			// RegisterShopOwner shop_owner - i registrasiya etmek ucin ulanylyar.
			shopOwner.POST("register", controllers.RegisterShopOwner)
			// admin.POST("register", middlewares.IsSuperAdmin(), adminController.RegisterAdmin)

			// LoginShopOwner shop_owner - i login etmek ucin ulanylyar.
			shopOwner.POST("login", controllers.LoginShopOwner)

			// UpdateShopOwner shop_owner - in maglumatlaryny uytgetmek ucin ulanylyar.
			shopOwner.PUT("update", controllers.UpdateShopOwner)

		}
	}

}
