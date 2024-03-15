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

	if err := models.ValidateDimension(dimension, false); err != nil {
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

func UpdateDimension(c *gin.Context) {
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

	if err := models.ValidateDimension(dimension, true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE dimensions SET dimension=$1, dimension_group_id=$2 WHERE id=$3", dimension.Dimension, dimension.DimensionGroupID, dimension.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func GetDimensionByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden id alynyar
	dimensionID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var dimension models.Dimension
	db.QueryRow(context.Background(), "SELECT id,dimension,dimension_group_id FROM dimensions WHERE id = $1 AND deleted_at IS NULL", dimensionID).Scan(&dimension.ID, &dimension.Dimension, &dimension.DimensionGroupID)

	// eger databse sol maglumat yok bolsa error return edilyar
	if dimension.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"dimension": dimension,
	})
}
