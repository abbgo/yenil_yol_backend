package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateComplaint(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var complaint models.Complaint
	if err := c.BindJSON(&complaint); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda complaints tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO complaints (text_tm,text_ru) VALUES ($1,$2)", complaint.TextTM, complaint.TextRU)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}

func UpdateComplaintByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var complaint models.Complaint
	if err := c.BindJSON(&complaint); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("complaints", complaint.ID, "NULL", db); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE complaints SET text_tm=$1 , text_ru=$2 WHERE id=$3", complaint.TextTM, complaint.TextRU, complaint.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}
