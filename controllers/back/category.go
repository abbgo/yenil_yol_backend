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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if category.Image == "" {
		image = nil
	} else {
		image = category.Image
	}

	// eger maglumatlar dogry bolsa onda categories tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	_, err = db.Exec(context.Background(), "INSERT INTO categories (name_tm,name_ru,image,slug_tm,slug_ru) VALUES ($1,$2,$3,$4,$5)", category.NameTM, category.NameRU, image, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa category ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var category models.CategoryUpdate
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var categoryID string
	var oldCategoryImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM categories WHERE id = $1 AND deleted_at IS NULL", category.ID).Scan(&categoryID, &oldCategoryImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if categoryID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
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
		if fileName != "" {
			// sonra helper_images tablisa category ucin gosulan surat pozulyar
			_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", fileName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}

			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + oldCategoryImage.String); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}
		}
		fileName = category.Image
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE categories SET name_tm=$1 , name_ru=$2 , image=$3 , slug_tm=$4 , slug_ru=$5 WHERE id=$6", category.NameTM, category.NameRU, fileName, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"), category.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}
