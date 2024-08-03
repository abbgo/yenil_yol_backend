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

func GetComplaintByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden brend id alynyar
	complaintID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var complaint models.Complaint
	db.QueryRow(context.Background(), "SELECT id,text_tm,text_ru FROM complaints WHERE id = $1 AND deleted_at IS NULL", complaintID).Scan(&complaint.ID, &complaint.TextTM, &complaint.TextRU)

	// eger databse sol maglumat yok bolsa error return edilyar
	if complaint.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"complaint": complaint,
	})
}

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

func DeletePermanentlyComplaintByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den brend id alynyar
	ID := c.Param("id")

	if err := helpers.ValidateRecordByID("complaints", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// brend - in suraty pozulandan sonra database - den brend pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM complaints WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}
