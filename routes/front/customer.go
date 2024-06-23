package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"
	"github/abbgo/yenil_yol/backend/middlewares"

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
		customerRoutes.PUT("update", middlewares.CheckTokenAdminOrCustomer(), controllers.UpdateCustomer)

		// UpdateCustomerPassword -> Customer - in maglumatlaryny uytgetmek ucin ulanylyar.
		customerRoutes.PUT("update-password", middlewares.CheckTokenAdminOrCustomer(), controllers.UpdateCustomerPassword)

		// // Customer - in access tokenin tazelelap refresh bilen access tokeni bile bermek
		// // ucin ulanylyar
		// customerRoutes.POST("refresh", helpers.RefreshTokenForAdmin)

		// GetCustomer -> haeder - den gelen id boyunca bir sany customer - i almak ucin ulanylyar.
		customerRoutes.GET(":id", middlewares.CheckTokenAdminOrCustomer(), controllers.GetCustomer)

		// GetCustomers -> hemme Customer - leri almak ucin ulanylyar.
		customerRoutes.GET("", middlewares.CheckToken("admin"), controllers.GetCustomers)

		// DeleteCustomerByID -> id boyunca Customer - i korzina salmak ucin ulanylyar.
		customerRoutes.DELETE(":id", middlewares.CheckToken("admin"), controllers.DeleteCustomerByID)

		// RestoreCustomerByID -> id boyunca customer - i korzinadan cykarmak ucin ulanylyar.
		customerRoutes.GET(":id/restore", middlewares.CheckToken("admin"), controllers.RestoreCustomerByID)

		// DeletePermanentlyCustomerByID -> id boyunca customer - i doly ( korzinadan ) pozmak ucin ulanylyar
		customerRoutes.DELETE(":id/delete", middlewares.CheckToken("admin"), controllers.DeletePermanentlyCustomerByID)
	}

}
