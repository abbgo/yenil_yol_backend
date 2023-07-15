package routes

import (
	// backApi "github/abbgo/isleg/isleg-backend/routes/back"
	adminApi "github/abbgo/yenil_yol/backend/routes/admin"
	backApi "github/abbgo/yenil_yol/backend/routes/back"

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
		adminApi.ShopOwnerRoutes(admin)

		// bu route - ler admin - ler ucin doredilen route - laryn toplumy
		adminApi.AdminRoutes(admin)
	}

	back := routes.Group("/api/back")
	{
		// bu route - ler magazynyn eyeleri ucin doredilen rout - laryn toplumy
		backApi.BackShopOwnerRoutes(back)

		// bu route - ler surat gosmak we pozmak ucin doredilen rout - laryn toplumy
		backApi.BackImageRoutes(back)

		// bu route - ler brend - ler ucin doredilen rout - laryn toplumy
		backApi.BrendRoutes(back)

		// bu route - ler category - ler ucin doredilen rout - laryn toplumy
		backApi.BrendCategories(back)
	}

	return routes

}
