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

func CreateBrend(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var brend models.Brend
	if err := c.BindJSON(&brend); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if brend.Image == "" {
		image = nil
	} else {
		image = brend.Image
	}

	// eger maglumatlar dogry bolsa onda brends tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO brends (name,image,slug) VALUES ($1,$2,$3)", brend.Name, image, slug.MakeLang(brend.Name, "en"))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa brend ucin gosulan surat pozulyar
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

func UpdateBrendByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var brend models.BrendUpdate
	if err := c.BindJSON(&brend); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var brendID string
	var olBrendImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM brends WHERE id = $1 AND deleted_at IS NULL", brend.ID).Scan(&brendID, &olBrendImage); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if brendID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// bu yerde brend - in suraty ucin fileName atly uytgeyan ululyk doredilyar
	// eger request body - dan surat gelmese onda brend - in suraty uytgedilmeyar diymek bolyar
	// sonun ucin brend - in onki suratyny goyyas , eger request body - dan surat gelen bolsa
	// onda taze suraty kone surat bilen calysyas
	var fileName string
	if brend.Image == "" {
		fileName = olBrendImage.String
	} else {
		// sonra helper_images tablisa brend ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", brend.Image)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if olBrendImage.String != "" {
			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + olBrendImage.String); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		fileName = brend.Image
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE brends SET name=$1 , image=$2 , slug=$3 WHERE id=$4", brend.Name, fileName, slug.MakeLang(brend.Name, "en"), brend.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}

func GetBrendByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden brend id alynyar
	brendID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var brend models.Brend
	var brendImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,name,image FROM brends WHERE id = $1 AND deleted_at IS NULL", brendID).Scan(&brend.ID, &brend.Name, &brendImage); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if brend.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	if brendImage.String != "" {
		brend.Image = brendImage.String
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"brend":  brend,
	})

}

func GetBrends(c *gin.Context) {

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

	// request query - den brend status alynyar
	// status -> brend pozulan ya-da pozulanmadygyny anlatyar
	// true bolsa pozulan
	// false bolsa pozulmadyk
	statusQuery := c.DefaultQuery("status", "false")
	status, err := strconv.ParseBool(statusQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora brend - leryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM brends WHERE deleted_at IS NULL`
	if status {
		queryCount = `SELECT COUNT(id) FROM brends WHERE deleted_at IS NOT NULL`
	}
	// database - den brend - laryn sany alynyar
	var countOfBrends uint
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfBrends); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora brend - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name,image FROM brends WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if status {
		rowQuery = `SELECT id,name,image FROM brends WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsBrend, err := db.Query(context.Background(), rowQuery, limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsBrend.Close()

	var brends []models.Brend
	for rowsBrend.Next() {
		var brend models.Brend
		var brendImage sql.NullString
		if err := rowsBrend.Scan(&brend.ID, &brend.Name, &brendImage); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if brendImage.String != "" {
			brend.Image = brendImage.String
		}
		brends = append(brends, brend)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"brens":  brends,
		"total":  countOfBrends,
	})

}

func DeleteBrendByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den brend id alynyar
	ID := c.Param("id")

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM brends WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa shop we sol brend - in deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE brends SET deleted_at = NOW() WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}

func RestoreBrendByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den brend id alynyar
	ID := c.Param("id")

	// alynan id den bolan brend database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM brends WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database sol id degisli brend yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa brend restore edilyar
	_, err = db.Exec(context.Background(), "UPDATE brends SET deleted_at = NULL WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})

}

func DeletePermanentlyBrendByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den brend id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	var image sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM brends WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id, &image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli brend yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger shop bar bolsa sonda brend - in suraty papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// brend - in suraty pozulandan sonra database - den brend pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM brends WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}
