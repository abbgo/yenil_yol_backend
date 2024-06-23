package adminApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/admin"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ShopOwnerRoutes(back *gin.RouterGroup) {
	shopOwner := back.Group("/shop-owners")
	{
		// RegisterShopOwner shop_owner - i registrasiya etmek ucin ulanylyar.
		shopOwner.POST("register", controllers.RegisterShopOwner)
		// admin.POST("register", middlewares.IsSuperAdmin(), adminController.RegisterAdmin)

		// LoginShopOwner shop_owner - i login etmek ucin ulanylyar.
		shopOwner.POST("login", controllers.LoginShopOwner)

		// UpdateShopOwner shop_owner - in maglumatlaryny uytgetmek ucin ulanylyar.
		shopOwner.PUT("update", middlewares.CheckTokenAdminOrShopOwner(), controllers.UpdateShopOwner)

		// UpdateShopOwnerPassword shop owner - in parolyny uytgetmek ucin ulanylyar.
		shopOwner.PUT("update-password", middlewares.CheckTokenAdminOrShopOwner(), controllers.UpdateShopOwnerPassword)

		// // ShopOwner - in access tokenin tazelelap refresh bilen access tokeni bile bermek
		// // ucin ulanylyar
		// shopOwner.POST("refresh", helpers.RefreshTokenForAdmin)

		// GetShopOwner funksiya haeder - den gelen id boyunca bir sany shop_owneri almak ucin ulanylyar.
		shopOwner.GET(":id", middlewares.CheckTokenAdminOrShopOwner(), controllers.GetShopOwner)

		// GetShopOwners funksiya hemme shop_owner - leri almak ucin ulanylyar.
		shopOwner.GET("", middlewares.CheckToken("admin"), controllers.GetShopOwners)

		// DeleteShopOwnerByID funksiya id boyunca shop_owner - i korzina salmak ucin ulanylyar.
		shopOwner.DELETE(":id", middlewares.CheckToken("admin"), controllers.DeleteShopOwnerByID)

		// RestoreShopOwnerByID funksiya id boyunca shop_owner - i korzinadan cykarmak ucin ulanylyar.
		shopOwner.GET(":id/restore", middlewares.CheckToken("admin"), controllers.RestoreShopOwnerByID)

		// DeletePermanentlyShopOwnerByID funksiya id boyunca shop_owner - i doly ( korzinadan ) pozmak ucin ulanylyar
		shopOwner.DELETE(":id/delete", middlewares.CheckToken("admin"), controllers.DeletePermanentlyShopOwnerByID)
	}
}
