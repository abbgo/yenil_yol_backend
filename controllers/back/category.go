package controllers

import (
	"context"
	"database/sql"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func CreateCategory(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCategory(category, true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda categories tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	_, err = db.Exec(context.Background(), "INSERT INTO categories (name_tm,name_ru,image,slug_tm,slug_ru,dimension_group_id,parent_category_id) VALUES ($1,$2,$3,$4,$5,$6,$7)", category.NameTM, category.NameRU, category.Image, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"), category.DimensionGroupID, category.ParentCategoryID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// category - nyn maglumatlary gosulandan sonra suraty bar bolsa helper_images tablisa category ucin gosulan surat pozulyar
	if category.Image.String != "" {
		if err := DeleteImageFromDB(category.Image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}

func UpdateCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCategory(category, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE categories SET name_tm=$1 , name_ru=$2 , image=$3 , slug_tm=$4 , slug_ru=$5, dimension_group_id=$6, parent_category_id=$7 WHERE id=$8", category.NameTM, category.NameRU, category.Image, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"), category.DimensionGroupID, category.ParentCategoryID, category.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func GetCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden category id alynyar
	categoryID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var category models.Category
	if err := db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,image,dimension_group_id,parent_category_id FROM categories WHERE id = $1 AND deleted_at IS NULL", categoryID).Scan(&category.ID, &category.NameTM, &category.NameRU, &category.Image, &category.DimensionGroupID, &category.ParentCategoryID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if category.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"category": category,
	})
}

func GetCategories(c *gin.Context) {
	var requestQuery helpers.StandartQuery
	var countOfCategories uint
	var categories []models.Category

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request query - den maglumatlar bind edilyar
	if err := c.Bind(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := requestQuery.Limit * (requestQuery.Page - 1)

	// request query - den status - a gora category - laryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL`
	if requestQuery.IsDeleted {
		queryCount = `SELECT COUNT(id) FROM categories WHERE deleted_at IS NOT NULL`
	}
	// database - den category - laryn sany alynyar
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfCategories); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora category - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name_tm,name_ru,image,dimension_group_id,parent_category_id FROM categories WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if requestQuery.IsDeleted {
		rowQuery = `SELECT id,name_tm,name_ru,image,dimension_group_id,parent_category_id FROM categories WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsCategory, err := db.Query(context.Background(), rowQuery, requestQuery.Limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()

	for rowsCategory.Next() {
		var category models.Category
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &category.NameRU, &category.Image, &category.DimensionGroupID, &category.ParentCategoryID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"categories": categories,
		"total":      countOfCategories,
	})
}

func DeleteCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den category id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("categories", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa categories we sol tablisa bilen baglanysykly tablisalaryn deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_category($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den category id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("categories", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa categories we sol tablisa bilen baglanysykly tablisalaryn deleted_at - ine NULL goyulyar
	_, err = db.Exec(context.Background(), "CALL restore_category($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den category id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	var image sql.NullString
	db.QueryRow(context.Background(), "SELECT id,image FROM categories WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id, &image)

	// eger database - de gelen id degisli category yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger category - nyn suraty bar bolsa onda ol local papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := os.Remove(helpers.ServerPath + "assets/" + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// category - nyn suraty pozulandan son category we onun bilen baglanysykly maglumatlar pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}
