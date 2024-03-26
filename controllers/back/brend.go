package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"

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

	// eger maglumatlar dogry bolsa onda brends tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO brends (name,image,slug) VALUES ($1,$2,$3)", brend.Name, brend.Image, slug.MakeLang(brend.Name, "en"))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa brend ucin gosulan surat pozulyar
	if brend.Image != "" {
		if err := DeleteImageFromDB(brend.Image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
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
	var brend models.Brend
	if err := c.BindJSON(&brend); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("brends", brend.ID, "NULL", db); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE brends SET name=$1 , image=$2 , slug=$3 WHERE id=$4", brend.Name, brend.Image, slug.MakeLang(brend.Name, "en"), brend.ID)
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
	db.QueryRow(context.Background(), "SELECT id,name,image FROM brends WHERE id = $1 AND deleted_at IS NULL", brendID).Scan(&brend.ID, &brend.Name, &brend.Image)

	// eger databse sol maglumat yok bolsa error return edilyar
	if brend.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"brend":  brend,
	})
}

func GetBrends(c *gin.Context) {
	var requestQuery helpers.StandartQuery
	var countOfBrends uint
	var brends []models.Brend

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

	// request query - den status - a gora brend - leryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM brends WHERE deleted_at IS NULL`
	if requestQuery.IsDeleted {
		queryCount = `SELECT COUNT(id) FROM brends WHERE deleted_at IS NOT NULL`
	}
	// database - den brend - laryn sany alynyar
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfBrends); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora brend - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name,image FROM brends WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if requestQuery.IsDeleted {
		rowQuery = `SELECT id,name,image FROM brends WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsBrend, err := db.Query(context.Background(), rowQuery, requestQuery.Limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsBrend.Close()

	for rowsBrend.Next() {
		var brend models.Brend
		if err := rowsBrend.Scan(&brend.ID, &brend.Name, &brend.Image); err != nil {
			helpers.HandleError(c, 400, err.Error())
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
	if err := helpers.ValidateRecordByID("brends", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
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
	if err := helpers.ValidateRecordByID("brends", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
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
	var brend models.Brend
	db.QueryRow(context.Background(), "SELECT id,image FROM brends WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&brend.ID, &brend.Image)

	// eger database - de gelen id degisli brend yok bolsa error return edilyar
	if brend.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger brend bar bolsa sonda brend - in suraty papkadan pozulyar
	if err := os.Remove(helpers.ServerPath + brend.Image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	if err := os.Remove(helpers.ServerPath + "assets/" + brend.Image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
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
