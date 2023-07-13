package controllers

import (
	"context"
	"database/sql"
	"fmt"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var brend models.Brend
	if err := c.BindJSON(&brend); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if brend.Image == "" {
		image = nil
	} else {
		image = brend.Image
	}

	// eger maglumatlar dogry bolsa onda brends tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	_, err = db.Exec(context.Background(), "INSERT INTO brends (name,image,slug) VALUES ($1,$2,$3)", brend.Name, image, slug.MakeLang(brend.Name, "en"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa shop ucin gosulan surat pozulyar
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

func UpdateBrendByID(c *gin.Context) {

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
	var brend models.BrendUpdate
	if err := c.BindJSON(&brend); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var brendID string
	var olBrendImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM brends WHERE id = $1 AND deleted_at IS NULL", brend.ID).Scan(&brendID, &olBrendImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if brendID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
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
		if fileName != "" {
			// sonra helper_images tablisa brend ucin gosulan surat pozulyar
			_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", fileName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}

			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + olBrendImage.String); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}
		}
		fileName = brend.Image
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE brends SET name=$1 , image=$2 , slug=$3 WHERE id=$4", brend.Name, fileName, slug.MakeLang(brend.Name, "en"), brend.ID)
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

func GetBrendByID(c *gin.Context) {

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

	// request parametrden brend id alynyar
	brendID := c.Param("id")
	fmt.Println("brendID: ", brendID)

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var brend models.Brend
	if err := db.QueryRow(context.Background(), "SELECT id,name,image FROM brends WHERE id = $1 AND deleted_at IS NULL", brendID).Scan(&brend.ID, &brend.Name, &brend.Image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "yalnys 1",
		})
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if brend.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request parametr - den limit alynyar
	limitStr := c.Query("limit")
	if limitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "limit is required",
		})
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// request parametr - den page alynyar
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "page is required",
		})
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// request query - den shop status alynyar
	// status -> shop pozulan ya-da pozulanmadygyny anlatyar
	// true bolsa pozulan
	// false bolsa pozulmadyk
	statusQuery := c.DefaultQuery("status", "false")
	status, err := strconv.ParseBool(statusQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// request query - den status - a gora brend - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name,image FROM brends WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if status {
		rowQuery = `SELECT id,name,image FROM brends WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsShop, err := db.Query(context.Background(), rowQuery, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "yalnys 2",
		})
		return
	}
	defer rowsShop.Close()

	var brends []models.Brend
	for rowsShop.Next() {
		var brend models.Brend
		if err := rowsShop.Scan(&brend.ID, &brend.Name, &brend.Image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
		brends = append(brends, brend)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"brens":  brends,
		"total":  countOfBrends,
	})

}
