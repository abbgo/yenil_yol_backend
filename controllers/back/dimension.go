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

func GetDimensionsByGroupID(c *gin.Context) {
	var dimensions []models.Dimension
	var dimensionQuery models.DimensionQuery

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request query - den maglumatlar helpers.StandartQuery struct boyunca bind edilyar
	if err := c.Bind(&dimensionQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&dimensionQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("dimension_groups", dimensionQuery.DimensionGroupID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// request query - den status - a gora dimension - lary almak ucin query yazylyar
	rowQuery := `SELECT id,dimension,dimension_group_id FROM dimensions WHERE deleted_at IS NULL AND dimension_group_id=$1`
	if dimensionQuery.IsDeleted {
		rowQuery = `SELECT id,dimension,dimension_group_id FROM dimensions WHERE deleted_at IS NOT NULL AND dimension_group_id=$1`
	}
	// database - den brend - lar alynyar
	rows, err := db.Query(context.Background(), rowQuery, dimensionQuery.DimensionGroupID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dimension models.Dimension
		if err := rows.Scan(&dimension.ID, &dimension.Dimension, &dimension.DimensionGroupID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		dimensions = append(dimensions, dimension)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"dimensions": dimensions,
	})
}

func DeleteDimensionByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("dimensions", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa dimension deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE dimensions SET deleted_at=NOW() WHERE id=$1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreDimensionByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("dimensions", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa dimension deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE dimensions SET deleted_at=NULL WHERE id=$1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}
