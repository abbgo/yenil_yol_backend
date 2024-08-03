package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetComplaints(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	var complaints []models.Complaint
	rows, err := db.Query(context.Background(), `SELECT id,text_tm,text_ru FROM complaints`)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var complaint models.Complaint
		if err := rows.Scan(&complaint.ID, &complaint.TextTM, &complaint.TextRU); err != nil {
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
