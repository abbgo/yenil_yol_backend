package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateComplaintProduct(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var cp models.ComplaintProduct
	if err := c.BindJSON(&cp); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("complaints", cp.ComplaintID, "NULL", db); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	if err := helpers.ValidateRecordByID("products", cp.ProductID, "NULL", db); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda complaints tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO complaint_products (complaint_id,product_id) VALUES ($1,$2)", cp.ComplaintID, cp.ProductID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}
