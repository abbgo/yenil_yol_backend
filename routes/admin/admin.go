package adminApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/admin"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(back *gin.RouterGroup) {

	admin := back.Group("/admins")
	{
		{
			// RegisterAdmin admin - i registrasiya etmek ucin ulanylyar.
			admin.POST("register", middlewares.IsSuperAdmin(), controllers.RegisterAdmin)
			// 	// admin.POST("register", middlewares.IsSuperAdmin(), adminController.RegisterAdmin)

			// LoginAdmin admin - i login etmek ucin ulanylyar.
			admin.POST("login", controllers.LoginAdmin)

			// UpdateAdmin admin - in maglumatlaryny uytgetmek ucin ulanylyar.
			admin.PUT("update", middlewares.IsSuperAdmin(), controllers.UpdateAdmin)

			// UpdateAdmin admin - in maglumatlaryny uytgetmek ucin ulanylyar.
			admin.PUT("update-password", middlewares.IsSuperAdmin(), controllers.UpdateAdminPassword)

			// Adminin - in access tokenin tazelelap refresh bilen access tokeni bile bermek
			// ucin ulanylyar
			admin.POST("refresh", middlewares.CheckAdmin(), helpers.RefreshTokenForAdmin)

			// GetAdmin funksiya haeder - den gelen id boyunca bir sany admin - i almak ucin ulanylyar.
			admin.GET("one", middlewares.CheckAdmin(), controllers.GetAdmin)

			// GetAdmins funksiya hemme admin - leri almak ucin ulanylyar.
			admin.GET("", middlewares.IsSuperAdmin(), controllers.GetAdmins)

			// DeleteAdminByID funksiya id boyunca admin - i korzina salmak ucin ulanylyar.
			admin.DELETE(":id", middlewares.IsSuperAdmin(), controllers.DeleteAdminByID)

			// RestoreAdminByID funksiya id boyunca admin - i korzinadan cykarmak ucin ulanylyar.
			admin.GET(":id/restore", middlewares.IsSuperAdmin(), controllers.RestoreAdminByID)

			// DeletePermanentlyAdminByID funksiya id boyunca admin - i doly ( korzinadan ) pozmak ucin ulanylyar
			admin.DELETE(":id/delete", middlewares.IsSuperAdmin(), controllers.DeletePermanentlyAdminByID)

		}
	}

}
