package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetComplaintProducts(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request param - dan shop owner id - ni alyarys
	shopOwnerID := c.Param("shopOwnerID")

	// shop owner id boyunca harydyn sikayatlarynyn db - dan alyarys
	var cps []serializations.ProductComplaint
	rows, err := db.Query(
		context.Background(),
		`SELECT DISTINCT ON (p.id) p.id,p.name_tm,p.name_ru FROM products p 
		INNER JOIN complaint_products cp ON cp.product_id=p.id  
		INNER JOIN shops s ON s.id=p.shop_id 
		WHERE s.shop_owner_id=$1 
		AND s.deleted_at IS NULL 
		AND p.deleted_at IS NULL 
		AND cp.deleted_at IS NULL`, shopOwnerID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cp serializations.ProductComplaint
		if err := rows.Scan(&cp.ProductID, &cp.NameTM, &cp.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// sikayat edilen harydyn suratyny alyarys
		db.QueryRow(
			context.Background(),
			`SELECT image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id WHERE pc.product_id=$1 AND pc.order_number=1 AND pi.order_number=1`,
			cp.ProductID).
			Scan(&cp.Image)

		// sikayat edilen haryda degisli sikayatlaryn sanyny alyarys
		db.QueryRow(context.Background(), `SELECT COUNT(id) FROM complaint_products WHERE product_id=$1`, cp.ProductID).Scan(&cp.ComplaintCount)

		cps = append(cps, cp)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             true,
		"complaint_products": cps,
	})
}
