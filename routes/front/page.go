package frontApi

import (
	controllers "github/abbgo/yenil_yol/backend/controllers/front"

	"github.com/gin-gonic/gin"
)

func PageRoutes(front *gin.RouterGroup) {

	pageRoutes := front.Group("/pages")
	{
		// GetPageByName -> request param - den gelyan ady boyunca
		// fronta sahypanyn terjimelerini ugratmak ucin
		// ulanylyar
		pageRoutes.GET(":name", controllers.GetPageByName)

	}

}
