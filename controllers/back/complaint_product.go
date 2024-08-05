package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetComplaintProduct(c *gin.Context) {
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
	var complaints []models.Complaint
	rows, err := db.Query(
		context.Background(),
		`SELECT c.text_tm,c.text_ru FROM complaints c 
		INNER JOIN complaint_products cp ON cp.complaint_id=c.id 
		INNER JOIN products p ON p.id=cp.product_id 
		INNER JOIN shops s ON s.id=p.shop_id 
		WHERE s.shop_owner_id=$1 
		AND c.deleted_at IS NULL 
		AND cp.deleted_at IS NULL 
		AND p.deleted_at IS NULL 
		AND s.deleted_at IS NULL 
		ORDER BY cp.created_at DESC`, shopOwnerID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var complaint models.Complaint
		if err := rows.Scan(&complaint.TextTM, &complaint.TextRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		complaints = append(complaints, complaint)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"complaints": complaints,
	})
}
