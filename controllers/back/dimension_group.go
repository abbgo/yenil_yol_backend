package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDimensionGroup(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var dimensionGroup models.DimensionGroup
	if err := c.BindJSON(&dimensionGroup); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda dimension_groups tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO dimension_groups (name) VALUES ($1)", dimensionGroup.Name)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}

func UpdateDimensionGroup(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var dimensionGroup models.DimensionGroup
	if err := c.BindJSON(&dimensionGroup); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("dimension_groups", dimensionGroup.ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE dimension_groups SET name=$1 WHERE id=$2", dimensionGroup.Name, dimensionGroup.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func GetDimensionGroupByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden brend id alynyar
	dimensionGroupID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var dimensionGroup models.DimensionGroup
	db.QueryRow(context.Background(), "SELECT id,name FROM dimension_groups WHERE id = $1 AND deleted_at IS NULL", dimensionGroupID).Scan(&dimensionGroup.ID, &dimensionGroup.Name)

	// eger databse sol maglumat yok bolsa error return edilyar
	if dimensionGroup.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          true,
		"dimension_group": dimensionGroup,
	})
}
