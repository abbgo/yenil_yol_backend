package shopOwnerApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/admin"
	"github/abbgo/yenil_yol/backend/helpers"

	"github.com/gin-gonic/gin"
)

func ShopOwnerRoutes(back *gin.RouterGroup) {

	shopOwner := back.Group("/shop-owners")
	{
		{
			// RegisterShopOwner shop_owner - i registrasiya etmek ucin ulanylyar.
			shopOwner.POST("register", controllers.RegisterShopOwner)
			// admin.POST("register", middlewares.IsSuperAdmin(), adminController.RegisterAdmin)

			// LoginShopOwner shop_owner - i login etmek ucin ulanylyar.
			shopOwner.POST("login", controllers.LoginShopOwner)

			// UpdateShopOwner shop_owner - in maglumatlaryny uytgetmek ucin ulanylyar.
			shopOwner.PUT("update", controllers.UpdateShopOwner)

			// ShopOwner - in access tokenin tazelelap refresh bilen access tokeni bile bermek
			// ucin ulanylyar
			shopOwner.POST("refresh", helpers.RefreshTokenForAdmin)

			// GetShopOwner funksiya haeder - den gelen id boyunca bir sany shop_owneri almak ucin ulanylyar.
			shopOwner.GET("one", controllers.GetShopOwner)

			// GetShopOwners funksiya hemme shop_owner - leri almak ucin ulanylyar.
			shopOwner.GET("", controllers.GetShopOwners)

			// DeleteShopOwnerByID funksiya id boyunca shop_owner - i korzina salmak ucin ulanylyar.
			shopOwner.DELETE(":id", controllers.DeleteShopOwnerByID)

			// RestoreShopOwnerByID funksiya id boyunca shop_owner - i korzinadan cykarmak ucin ulanylyar.
			shopOwner.GET(":id/restore", controllers.RestoreShopOwnerByID)

			// DeletePermanentlyShopOwnerByID funksiya id boyunca shop_owner - i duybinden pozmak ucin ulanylyar
			shopOwner.DELETE(":id/delete", controllers.DeletePermanentlyShopOwnerByID)

		}
	}

}
