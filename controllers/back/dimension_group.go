package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"github/abbgo/yenil_yol/backend/serializations"
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

func GetDimensionGroupsWithDimensions(c *gin.Context) {
	var dimensionGroups []models.DimensionGroup
	var requestQuery helpers.StandartQuery
	var IsDeleted = "NULL"

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request query - den maglumatlar helpers.StandartQuery struct boyunca bind edilyar
	if err := c.Bind(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if requestQuery.IsDeleted {
		IsDeleted = "NOT NULL"
	}

	// query - den gelyan limit we page boyunca databasede ulanyljak offset hasaplanyar
	offset := requestQuery.Limit * (requestQuery.Page - 1)

	// request query - den status - a gora dimension_group - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(`SELECT id,name FROM dimension_groups WHERE deleted_at IS %s ORDER BY created_at DESC LIMIT %v OFFSET %v`, IsDeleted, requestQuery.Limit, offset)

	// database - den brend - lar alynyar
	rows, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dimensionGroup models.DimensionGroup
		if err := rows.Scan(&dimensionGroup.ID, &dimensionGroup.Name); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		rowsDimensions, err := db.Query(context.Background(), `SELECT id,dimension FROM dimensions WHERE dimension_group_id=$1 AND deleted_at IS NULL`, dimensionGroup.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsDimensions.Close()

		for rowsDimensions.Next() {
			var dimension models.Dimension
			if err := rowsDimensions.Scan(&dimension.ID, &dimension.Dimension); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			dimensionGroup.Dimensions = append(dimensionGroup.Dimensions, dimension)
		}
		dimensionGroups = append(dimensionGroups, dimensionGroup)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           true,
		"dimension_groups": dimensionGroups,
	})
}

func GetDimensionGroupsWithDimensionsList(c *gin.Context) {
	var dimensionGroups []serializations.DimensionGroup
	var requestQuery helpers.StandartQuery

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request query - den maglumatlar helpers.StandartQuery struct boyunca bind edilyar
	if err := c.Bind(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// query - den gelyan limit we page boyunca databasede ulanyljak offset hasaplanyar
	offset := requestQuery.Limit * (requestQuery.Page - 1)

	// request query - den status - a gora dimension_group - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(
		`SELECT id,name FROM dimension_groups WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT %v OFFSET %v`,
		requestQuery.Limit, offset,
	)

	// database - den brend - lar alynyar
	rows, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dimensionGroup serializations.DimensionGroup
		if err := rows.Scan(&dimensionGroup.ID, &dimensionGroup.Name); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		rowsDimensions, err := db.Query(
			context.Background(), `SELECT dimension FROM dimensions WHERE dimension_group_id=$1 AND deleted_at IS NULL`,
			dimensionGroup.ID,
		)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsDimensions.Close()

		for rowsDimensions.Next() {
			var dimension string
			if err := rowsDimensions.Scan(&dimension); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			dimensionGroup.Dimensions = append(dimensionGroup.Dimensions, dimension)
		}
		dimensionGroups = append(dimensionGroups, dimensionGroup)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           true,
		"dimension_groups": dimensionGroups,
	})
}

func DeleteDimensionGroupByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den dimension_group id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("dimension_groups", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa dimension_group we sona degisli dimension - laryn deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_dimension_group($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreDimensionGroupByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den dimension_group id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("dimension_groups", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa dimension_group we sona degisli dimension - laryn deleted_at - ine NULL goyulyar
	_, err = db.Exec(context.Background(), "CALL restore_dimension_group($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyDimensionGroupByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den dimension_group id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("dimension_groups", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa dimension_group we sona degisli dimension - laryn hemmesi doly pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM dimension_groups WHERE id=$1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}
