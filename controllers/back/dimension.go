package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDimension(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var dimension models.Dimension
	if err := c.BindJSON(&dimension); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateDimension(dimension.DimensionGroupID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda dimensions tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO dimensions (dimension,dimension_group_id) VALUES ($1,$2)", dimension.Dimension, dimension.DimensionGroupID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}
