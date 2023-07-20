package frontApi

import (
	// backController "github/abbgo/isleg/isleg-backend/controllers/back"

	frontController "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(front *gin.RouterGroup) {

	customerRoutes := front.Group("/customers")
	{
		// RegisterCustomer -> kliendi registrasiya etmek ucin ulanylyar
		customerRoutes.POST("register", frontController.RegisterCustomer)

		// // LoginAdmin admin - i login etmek ucin ulanylyar.
		// customerRoutes.POST("login", controllers.LoginAdmin)

		// // UpdateAdmin admin - in maglumatlaryny uytgetmek ucin ulanylyar.
		// customerRoutes.PUT("update", middlewares.IsSuperAdmin(), controllers.UpdateAdmin)

		// // UpdateAdmin admin - in maglumatlaryny uytgetmek ucin ulanylyar.
		// customerRoutes.PUT("update-password", middlewares.IsSuperAdmin(), controllers.UpdateAdminPassword)

		// // Adminin - in access tokenin tazelelap refresh bilen access tokeni bile bermek
		// // ucin ulanylyar
		// customerRoutes.POST("refresh", middlewares.CheckAdmin(), helpers.RefreshTokenForAdmin)

		// // GetAdmin funksiya haeder - den gelen id boyunca bir sany shop_owneri almak ucin ulanylyar.
		// customerRoutes.GET("one", middlewares.CheckAdmin(), controllers.GetAdmin)

		// // GetAdmins funksiya hemme admin - leri almak ucin ulanylyar.
		// customerRoutes.GET("", middlewares.IsSuperAdmin(), controllers.GetAdmins)

		// // DeleteAdminByID funksiya id boyunca admin - i korzina salmak ucin ulanylyar.
		// customerRoutes.DELETE(":id", middlewares.IsSuperAdmin(), controllers.DeleteAdminByID)

		// // RestoreAdminByID funksiya id boyunca admin - i korzinadan cykarmak ucin ulanylyar.
		// customerRoutes.GET(":id/restore", middlewares.IsSuperAdmin(), controllers.RestoreAdminByID)

		// // DeletePermanentlyAdminByID funksiya id boyunca admin - i doly ( korzinadan ) pozmak ucin ulanylyar
		// customerRoutes.DELETE(":id/delete", middlewares.IsSuperAdmin(), controllers.DeletePermanentlyAdminByID)
		// SecuredCustomerRoutes(frontRoutes)
	}

}
