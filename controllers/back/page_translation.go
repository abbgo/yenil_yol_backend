package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePageTr(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var pageTr models.PageTranslation
	if err := c.BindJSON(&pageTr); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	var title_tm interface{}
	if pageTr.TitleTM == "" {
		title_tm = nil
	} else {
		title_tm = pageTr.TitleTM
	}

	var title_ru interface{}
	if pageTr.TitleRU == "" {
		title_ru = nil
	} else {
		title_ru = pageTr.TitleRU
	}

	var description_tm interface{}
	if pageTr.DescriptionTM == "" {
		description_tm = nil
	} else {
		description_tm = pageTr.DescriptionTM
	}

	var description_ru interface{}
	if pageTr.DescriptionRU == "" {
		description_ru = nil
	} else {
		description_ru = pageTr.DescriptionRU
	}

	// eger maglumatlar dogry bolsa onda page_translations tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO page_translations (title_tm,title_ru,description_tm,description_ru,page_id) VALUES ($1,$2,$3,$4,$5)", title_tm, title_ru, description_tm, description_ru, pageTr.PageID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// // brend - yn maglumatlary gosulandan sonra helper_images tablisa page ucin gosulan surat pozulyar
	// _, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", image)
	// if err != nil {
	// 	helpers.HandleError(c, 400, err.Error())
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})

}

func UpdatePageTrByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var pageTr models.PageTranslationUpdate
	if err := c.BindJSON(&pageTr); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var pageTrID string
	if err := db.QueryRow(context.Background(), "SELECT id FROM page_translations WHERE id = $1 AND deleted_at IS NULL", pageTr.ID).Scan(&pageTrID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if pageTrID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE page_translations SET title_tm=$1 , title_ru=$2 , description_tm=$3 , description_ru=$4 , page_id=$5 WHERE id=$6", pageTr.TitleTM, pageTr.TitleRU, pageTr.DescriptionTM, pageTr.DescriptionRU, pageTr.PageID, pageTr.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}
