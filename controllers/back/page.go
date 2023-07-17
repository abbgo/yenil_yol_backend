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

func GetPages(c *gin.Context) {

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

	// request query - den page status alynyar
	// status -> page pozulan ya-da pozulanmadygyny anlatyar
	// true bolsa pozulan
	// false bolsa pozulmadyk
	statusQuery := c.DefaultQuery("status", "false")
	status, err := strconv.ParseBool(statusQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora page - leryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM pages WHERE deleted_at IS NULL`
	if status {
		queryCount = `SELECT COUNT(id) FROM brends WHERE deleted_at IS NOT NULL`
	}
	// database - den page - laryn sany alynyar
	var countOfPages uint
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfPages); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora page - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name,image FROM pages WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if status {
		rowQuery = `SELECT id,name,image FROM pages WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsPage, err := db.Query(context.Background(), rowQuery, limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsPage.Close()

	var pages []models.Page
	for rowsPage.Next() {
		var page models.Page
		var pageImage sql.NullString
		if err := rowsPage.Scan(&page.ID, &page.Name, &pageImage); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if pageImage.String != "" {
			page.Image = pageImage.String
		}
		pages = append(pages, page)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pages":  pages,
		"total":  countOfPages,
	})

}

func DeletePageByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den page id alynyar
	ID := c.Param("id")

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM pages WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa shop we sol page - in deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE pages SET deleted_at = NOW() WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}

func RestorePageByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den page id alynyar
	ID := c.Param("id")

	// alynan id den bolan page database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM pages WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database sol id degisli brend yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa brend restore edilyar
	_, err = db.Exec(context.Background(), "UPDATE pages SET deleted_at = NULL WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})

}

func DeletePermanentlyPageByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den page id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	var image sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM pages WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id, &image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli page yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger page bar bolsa sonda page - in suraty papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// page - in suraty pozulandan sonra database - den page pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM pages WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}
