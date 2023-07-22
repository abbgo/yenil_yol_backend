package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strconv"

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

	// var title_tm interface{}
	// if pageTr.TitleTM == "" {
	// 	title_tm = nil
	// } else {
	// 	title_tm = pageTr.TitleTM
	// }

	// var title_ru interface{}
	// if pageTr.TitleRU == "" {
	// 	title_ru = nil
	// } else {
	// 	title_ru = pageTr.TitleRU
	// }

	// var text_title_tm interface{}
	// if pageTr.TextTitleTM == "" {
	// 	text_title_tm = nil
	// } else {
	// 	text_title_tm = pageTr.TextTitleTM
	// }

	// var text_title_ru interface{}
	// if pageTr.TextTitleRU == "" {
	// 	text_title_ru = nil
	// } else {
	// 	text_title_ru = pageTr.TextTitleRU
	// }

	// var description_tm interface{}
	// if pageTr.DescriptionTM == "" {
	// 	description_tm = nil
	// } else {
	// 	description_tm = pageTr.DescriptionTM
	// }

	// var description_ru interface{}
	// if pageTr.DescriptionRU == "" {
	// 	description_ru = nil
	// } else {
	// 	description_ru = pageTr.DescriptionRU
	// }

	// eger maglumatlar dogry bolsa onda page_translations tablisa maglumatlar gosulyar
	// _, err = db.Exec(context.Background(), "INSERT INTO page_translations (title_tm,title_ru,text_title_tm,text_title_ru,description_tm,description_ru,page_id) VALUES ($1,$2,$3,$4,$5)", title_tm, title_ru, text_title_tm, text_title_ru, description_tm, description_ru, pageTr.PageID)
	// if err != nil {
	// 	helpers.HandleError(c, 400, err.Error())
	// 	return
	// }

	_, err = db.Exec(context.Background(), "INSERT INTO page_translations (text_title_tm,text_title_ru,description_tm,description_ru,page_id,order_number) VALUES ($1,$2,$3,$4,$5,$6)", pageTr.TextTitleTM, pageTr.TextTitleRU, pageTr.DescriptionTM, pageTr.DescriptionRU, pageTr.PageID, pageTr.OrderNumber)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

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
	_, err = db.Exec(context.Background(), "UPDATE page_translations SET description_tm=$1 , description_ru=$2 , page_id=$3 , text_title_tm=$4 , text_title_ru=$5 WHERE id=$6", pageTr.DescriptionTM, pageTr.DescriptionRU, pageTr.PageID, pageTr.TextTitleTM, pageTr.TextTitleRU, pageTr.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}

func GetPageTrByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden page_translation id alynyar
	pageTrID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var pageTr models.PageTranslation
	if err := db.QueryRow(context.Background(), "SELECT id,text_title_tm,text_title_ru,description_tm,description_ru,page_id FROM page_translations WHERE id = $1 AND deleted_at IS NULL", pageTrID).Scan(
		&pageTr.ID,
		&pageTr.TextTitleTM,
		&pageTr.TextTitleRU,
		&pageTr.DescriptionTM,
		&pageTr.DescriptionRU,
		&pageTr.PageID,
	); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if pageTr.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           true,
		"page_translation": pageTr,
	})

}

func GetPageTrs(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// // request parametr - den page_id alynyar
	pageID := c.Query("page_id")
	if pageID == "" {
		helpers.HandleError(c, 400, "page_id is required")
		return
	}

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

	// request query - den status - a gora page - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf("SELECT id,text_title_tm,text_title_ru,description_tm,description_ru,page_id FROM page_translations WHERE deleted_at IS NULL AND page_id = %v", pageID)
	if status {
		rowQuery = fmt.Sprintf("SELECT id,text_title_tm,text_title_ru,description_tm,description_ru,page_id FROM page_translations WHERE deleted_at IS NOT NULL AND page_id = %v", pageID)
	}

	// database - den page_translation - lar alynyar
	rowsPageTr, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsPageTr.Close()

	var pageTrs []models.PageTranslation
	for rowsPageTr.Next() {
		var pageTr models.PageTranslation
		if err := rowsPageTr.Scan(&pageTr.ID, &pageTr.TextTitleTM, &pageTr.TextTitleRU, &pageTr.DescriptionTM, &pageTr.DescriptionRU, &pageTr.PageID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		pageTrs = append(pageTrs, pageTr)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":            true,
		"page_translations": pageTrs,
	})

}

func DeletePageTrByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den page_translation id alynyar
	ID := c.Param("id")

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM page_translations WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa shop we sol page_translation - in deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE page_translations SET deleted_at = NOW() WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}

func RestorePageTrByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den page_translation id alynyar
	ID := c.Param("id")

	// alynan id den bolan page database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM page_translations WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database sol id degisli page_translation yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa page_translation restore edilyar
	_, err = db.Exec(context.Background(), "UPDATE page_translations SET deleted_at = NULL WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})

}

func DeletePermanentlyPageTrByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den page_translation id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM page_translations WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli page_translation yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// page - in suraty pozulandan sonra database - den page_translation pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM page_translations WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}
