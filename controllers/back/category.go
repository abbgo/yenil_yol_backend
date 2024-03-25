package controllers

import (
	"context"
	"database/sql"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"
	"strconv"

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
	_, err = db.Exec(context.Background(), "INSERT INTO categories (name_tm,name_ru,image,slug_tm,slug_ru,dimension_group_id) VALUES ($1,$2,$3,$4,$5,$6)", category.NameTM, category.NameRU, category.Image, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"), category.DimensionGroupID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// category - nyn maglumatlary gosulandan sonra suraty bar bolsa helper_images tablisa category ucin gosulan surat pozulyar
	if category.Image != "" {
		if err := DeleteImageFromDB(category.Image); err != nil {
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

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var categoryID string
	var oldCategoryImage sql.NullString
	db.QueryRow(context.Background(), "SELECT id,image FROM categories WHERE id = $1 AND deleted_at IS NULL", category.ID).Scan(&categoryID, &oldCategoryImage)

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if categoryID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	if err := models.ValidateCategory(category, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// bu yerde category - in suraty ucin fileName atly uytgeyan ululyk doredilyar
	// eger request body - dan surat gelmese onda category - in suraty uytgedilmeyar diymek bolyar
	// sonun ucin category - in onki suratyny goyyas , eger request body - dan surat gelen bolsa
	// onda taze suraty kone surat bilen calysyas
	var fileName string
	if category.Image == "" {
		fileName = oldCategoryImage.String
	} else {
		// sonra helper_images tablisa category ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", category.Image)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if oldCategoryImage.String != "" {
			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + oldCategoryImage.String); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		fileName = category.Image
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE categories SET name_tm=$1 , name_ru=$2 , image=$3 , slug_tm=$4 , slug_ru=$5, dimension_group_id=$6 WHERE id=$7", category.NameTM, category.NameRU, fileName, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"), category.DimensionGroupID, category.ID)
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
	var categoryImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,name_tm,image FROM categories WHERE id = $1 AND deleted_at IS NULL", categoryID).Scan(&category.ID, &category.NameTM, &categoryImage); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if category.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	if categoryImage.String != "" {
		category.Image = categoryImage.String
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"category": category,
	})
}

func GetCategories(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den limit alynyar
	limitStr := c.Query("limit")
	if limitStr == "" {
		helpers.HandleError(c, 400, "limit is required")
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request parametr - den page alynyar
	pageStr := c.Query("page")
	if pageStr == "" {
		helpers.HandleError(c, 400, "page is required")
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// request query - den category status alynyar
	// status -> category pozulan ya-da pozulanmadygyny anlatyar
	// true bolsa pozulan
	// false bolsa pozulmadyk
	statusQuery := c.DefaultQuery("status", "false")
	status, err := strconv.ParseBool(statusQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora category - laryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL`
	if status {
		queryCount = `SELECT COUNT(id) FROM categories WHERE deleted_at IS NOT NULL`
	}
	// database - den category - laryn sany alynyar
	var countOfCategories uint
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfCategories); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora category - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name_tm,image FROM categories WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if status {
		rowQuery = `SELECT id,name_tm,image FROM categories WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsCategory, err := db.Query(context.Background(), rowQuery, limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()

	var categories []models.Category
	for rowsCategory.Next() {
		var category models.Category
		var categoryImage sql.NullString
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &categoryImage); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if categoryImage.String != "" {
			category.Image = categoryImage.String
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

	// hemme zat dogry bolsa shop we sol category - nin deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE categories SET deleted_at = NOW() WHERE id = $1", ID)
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

	// hemme zat dogry bolsa category restore edilyar
	_, err = db.Exec(context.Background(), "UPDATE categories SET deleted_at = NULL WHERE id = $1", ID)
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

	// eger shop bar bolsa sonda category - nin suraty papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// brend - in suraty pozulandan sonra database - den category pozulyar
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
