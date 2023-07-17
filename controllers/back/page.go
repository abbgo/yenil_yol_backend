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
)

func CreatePage(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var page models.Page
	if err := c.BindJSON(&page); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if page.Image == "" {
		image = nil
	} else {
		image = page.Image
	}

	// eger maglumatlar dogry bolsa onda pages tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO pages (name,image) VALUES ($1,$2)", page.Name, image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa page ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})

}

func UpdatePageByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var page models.PageUpdate
	if err := c.BindJSON(&page); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var pageID string
	var oldPageImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM pages WHERE id = $1 AND deleted_at IS NULL", page.ID).Scan(&pageID, &oldPageImage); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if pageID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// bu yerde page - in suraty ucin fileName atly uytgeyan ululyk doredilyar
	// eger request body - dan surat gelmese onda page - in suraty uytgedilmeyar diymek bolyar
	// sonun ucin page - in onki suratyny goyyas , eger request body - dan surat gelen bolsa
	// onda taze suraty kone surat bilen calysyas
	var fileName string
	if page.Image == "" {
		fileName = oldPageImage.String
	} else {
		// sonra helper_images tablisa page ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", page.Image)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if oldPageImage.String != "" {
			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + oldPageImage.String); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		fileName = page.Image
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE pages SET name=$1 , image=$2 WHERE id=$3", page.Name, fileName, page.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}

func GetPageByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden page id alynyar
	pageID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var page models.Page
	var pageImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,name,image FROM pages WHERE id = $1 AND deleted_at IS NULL", pageID).Scan(&page.ID, &page.Name, &pageImage); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if page.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	if pageImage.String != "" {
		page.Image = pageImage.String
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"page":   page,
	})

}
