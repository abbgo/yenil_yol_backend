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
		// bu route - ler musderiler ucin doredilen rout - laryn toplumy
		frontApi.CustomerRoutes(front)

		// bu route - ler halanlarym ucin doredilen rout - laryn toplumy
		frontApi.LikeRoutes(front)

		// bu route - ler shop - lar ucin doredilen rout - laryn toplumy
		frontApi.ShopRoutes(front)

		// bu route - ler fronta product maglumatlary ugratmak ucin doredilen rout - laryn toplumy
		frontApi.ProductRoutes(front)

		// bu route - ler fronta category maglumatlary ugratmak ucin doredilen rout - laryn toplumy
		frontApi.CategoryRoutes(front)

		// bu route - ler fronta sikayatlary ucin doredilen rout - laryn toplumy
		frontApi.ComplaintRoutes(front)

		// bu route - ler fronta harydyn sikayatlary ucin doredilen rout - laryn toplumy
		frontApi.ComplaintProductRoutes(front)
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
		backApi.CategoryRoutes(back)

		// bu route - ler saytyn sazlamalary ucin doredilen rout - laryn toplumy
		backApi.SettingRoutes(back)

		// bu route - ler product ucin doredilen rout - laryn toplumy
		backApi.ProductRoutes(back)

		// bu route - ler dimension groups ucin doredilen rout - laryn toplumy
		backApi.DimensionGroupRoutes(back)

		// bu route - ler dimension ucin doredilen rout - laryn toplumy
		backApi.DimensionRoutes(back)

		// bu route - ler sikayatlaryn teksti ucin doredilen rout - laryn toplumy
		backApi.ComplaintRoutes(back)

		// bu route - ler harydyn sikayatlary ucin doredilen rout - laryn toplumy
		backApi.ComplaintProductRoutes(back)
	}

	return routes
}
