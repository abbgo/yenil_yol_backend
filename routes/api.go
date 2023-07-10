package routes

import (
	// backApi "github/abbgo/isleg/isleg-backend/routes/back"
	shopOwnerApi "github/abbgo/yenil_yol/backend/routes/admin"
	frontApi "github/abbgo/yenil_yol/backend/routes/front"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {

	routes := gin.Default()

	// cors
	// routes.Use(cors.Default())

	routes.Use(gzip.Gzip(gzip.DefaultCompression))

	routes.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "RefreshToken", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		// MaxAge:           12 * time.Hour,
	}))

	front := routes.Group("/api")
	{
		// customer routes
		// frontApi.CustomerRoutes(front)

		// bu group - a degisli api - lerden maglumat alynanda ( :lang ) paramter boyunca uytgedilip
		// terjime alynyar
		frontApi.FrontRoutes(front)
	}

	admin := routes.Group("/api")
	{
		// bu rout - ler magazynyn eyeleri ucin doredilen rout - laryn toplumy
		shopOwnerApi.ShopOwnerRoutes(admin)
	}

	back := routes.Group("/api/back")
	{
		// bu rout - ler magazynyn eyeleri ucin doredilen rout - laryn toplumy
		shopOwnerApi.ShopOwnerRoutes(back)
	}

	return routes

}
