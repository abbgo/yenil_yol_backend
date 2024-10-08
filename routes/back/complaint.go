package back

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/back"
	"github/abbgo/yenil_yol/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ComplaintRoutes(back *gin.RouterGroup) {
	backComplaintApi := back.Group("/complaints").Use(middlewares.CheckToken("admin"))
	{
		// CreateComplaint -> sikayatyn tekstini gosmak ulanylar
		backComplaintApi.POST("", controllers.CreateComplaint)

		// UpdateComplaintByID -> id boyunca sikayatyn tekstini update etmek ucin ulanylyar
		backComplaintApi.PUT("", controllers.UpdateComplaintByID)

		// GetComplaintByID -> id - si boyunca sikayatyn tekstini almak ucin ulanylyar
		backComplaintApi.GET(":id", controllers.GetComplaintByID)

		// GetComplaints -> Ahli sikayatyn tekstini almak ucin ulanylyar
		backComplaintApi.GET("", controllers.GetComplaints)

		// DeletePermanentlyComplaintByID -> id boyunca sikayaty - i doly pozmak ucin ulanylyar
		backComplaintApi.DELETE(":id/delete", controllers.DeletePermanentlyComplaintByID)
	}
}
