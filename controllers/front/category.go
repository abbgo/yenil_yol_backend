package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategoriesShopID(c *gin.Context) {
	var categories []models.Category

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	shopID := c.Param("shop_id")

	// shop - a degisli category - ler alynyar
	rowsCategory, err := db.Query(context.Background(),
		`SELECT c.id,c.name_tm,c.name_ru FROM categories c 
		INNER JOIN shop_categories sc ON sc.category_id=c.id 
		WHERE sc.shop_id=$1 AND c.parent_category_id IS NULL AND c.deleted_at IS NULL AND sc.deleted_at IS NULL`,
		shopID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()
	for rowsCategory.Next() {
		var category models.Category
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &category.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"categories": categories,
	})
}
