package adminApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/admin"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(back *gin.RouterGroup) {

	admin := back.Group("/admins")
	{
		{
			// RegisterAdmin admin - i registrasiya etmek ucin ulanylyar.
			admin.POST("register", controllers.RegisterAdmin)
			// 	// admin.POST("register", middlewares.IsSuperAdmin(), adminController.RegisterAdmin)

			// LoginAdmin admin - i login etmek ucin ulanylyar.
			admin.POST("login", controllers.LoginAdmin)

			// 	// UpdateShopOwner shop_owner - in maglumatlaryny uytgetmek ucin ulanylyar.
			// 	shopOwner.PUT("update", controllers.UpdateShopOwner)

			// 	// ShopOwner - in access tokenin tazelelap refresh bilen access tokeni bile bermek
			// 	// ucin ulanylyar
			// 	shopOwner.POST("refresh", helpers.RefreshTokenForAdmin)

			// 	// GetShopOwner funksiya haeder - den gelen id boyunca bir sany shop_owneri almak ucin ulanylyar.
			// 	shopOwner.GET("one", controllers.GetShopOwner)

			// 	// GetShopOwners funksiya hemme shop_owner - leri almak ucin ulanylyar.
			// 	shopOwner.GET("", controllers.GetShopOwners)

			// 	// DeleteShopOwnerByID funksiya id boyunca shop_owner - i korzina salmak ucin ulanylyar.
			// 	shopOwner.DELETE(":id", controllers.DeleteShopOwnerByID)

			// 	// RestoreShopOwnerByID funksiya id boyunca shop_owner - i korzinadan cykarmak ucin ulanylyar.
			// 	shopOwner.GET(":id/restore", controllers.RestoreShopOwnerByID)

			// 	// DeletePermanentlyShopOwnerByID funksiya id boyunca shop_owner - i doly ( korzinadan ) pozmak ucin ulanylyar
			// 	shopOwner.DELETE(":id/delete", controllers.DeletePermanentlyShopOwnerByID)

		}
	}

}
