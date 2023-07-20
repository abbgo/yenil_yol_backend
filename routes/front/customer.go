package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(front *gin.RouterGroup) {

	customerRoutes := front.Group("/customers")
	{
		// RegisterCustomer -> kliendi registrasiya etmek ucin ulanylyar
		customerRoutes.POST("register", controllers.RegisterCustomer)

		// LoginCustomer -> Customer - i login etmek ucin ulanylyar.
		customerRoutes.POST("login", controllers.LoginCustomer)

		// UpdateCustomer -> Customer - in maglumatlaryny uytgetmek ucin ulanylyar.
		customerRoutes.PUT("update", controllers.UpdateCustomer)

		// // UpdateAdmin admin - in maglumatlaryny uytgetmek ucin ulanylyar.
		// customerRoutes.PUT("update-password",  controllers.UpdateAdminPassword)

		// // Adminin - in access tokenin tazelelap refresh bilen access tokeni bile bermek
		// // ucin ulanylyar
		// customerRoutes.POST("refresh", middlewares.CheckAdmin(), helpers.RefreshTokenForAdmin)

		// // GetAdmin funksiya haeder - den gelen id boyunca bir sany shop_owneri almak ucin ulanylyar.
		// customerRoutes.GET("one", middlewares.CheckAdmin(), controllers.GetAdmin)

		// // GetAdmins funksiya hemme admin - leri almak ucin ulanylyar.
		// customerRoutes.GET("",  controllers.GetAdmins)

		// // DeleteAdminByID funksiya id boyunca admin - i korzina salmak ucin ulanylyar.
		// customerRoutes.DELETE(":id",  controllers.DeleteAdminByID)

		// // RestoreAdminByID funksiya id boyunca admin - i korzinadan cykarmak ucin ulanylyar.
		// customerRoutes.GET(":id/restore",  controllers.RestoreAdminByID)

		// // DeletePermanentlyAdminByID funksiya id boyunca admin - i doly ( korzinadan ) pozmak ucin ulanylyar
		// customerRoutes.DELETE(":id/delete",  controllers.DeletePermanentlyAdminByID)
		// SecuredCustomerRoutes(frontRoutes)
	}

}
