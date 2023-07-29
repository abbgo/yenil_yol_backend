package controllers

import (
	"context"
	"database/sql"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponsePage struct {
	Image            string            `json:"image,omitempty"`
	TitleTM          string            `json:"title_tm,omitempty"`
	TitleRU          string            `json:"title_ru,omitempty"`
	PageTranslations []PageTranslation `json:"page_translations,omitempty"`
}

type PageTranslation struct {
	TextTitleTM   string `json:"text_title_tm,omitempty"`
	TextTitleRU   string `json:"text_title_ru,omitempty"`
	DescriptionTM string `json:"description_tm,omitempty"`
	DescriptionRU string `json:"description_ru,omitempty"`
}

func GetPageByName(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden page name alynyar
	pageName := c.Param("name")

	// database - den request parametr - den gelen name boyunca maglumat cekilyar
	var page ResponsePage
	var pageImage sql.NullString
	rowQuery := "SELECT p.image,p.title_tm,p.title_ru,pt.text_title_tm,pt.text_title_ru,pt.description_tm,pt.description_ru FROM pages p INNER JOIN page_translations pt ON pt.page_id = p.id WHERE p.deleted_at IS NULL AND pt.deleted_at IS NULL AND p.name = $1 ORDER BY pt.order_number ASC"
	rows, err := db.Query(context.Background(), rowQuery, pageName)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	for rows.Next() {
		var pageTranslation PageTranslation
		if err := rows.Scan(&pageImage, &page.TitleTM, &page.TitleRU, &pageTranslation.TextTitleTM, &pageTranslation.TextTitleRU, &pageTranslation.DescriptionTM, &pageTranslation.DescriptionRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if pageImage.String != "" {
			page.Image = pageImage.String
		}
		page.PageTranslations = append(page.PageTranslations, pageTranslation)
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if page.TitleTM == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"page":   page,
	})

}
