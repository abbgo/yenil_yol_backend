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
		`SELECT DISTINCT ON (c.id) c.id,c.name_tm,c.name_ru FROM categories c
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

		// child category alynyar
		queryForChildCategory := `SELECT DISTINCT ON (c.id) c.id,c.name_tm,c.name_ru,c.parent_category_id FROM categories c 
		INNER JOIN shop_categories sc ON sc.category_id=c.id 
		WHERE sc.shop_id=$1 AND c.parent_category_id=$2 AND c.deleted_at IS NULL AND sc.deleted_at IS NULL`

		rowsChildCategory, err := db.Query(context.Background(), queryForChildCategory, shopID, category.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsChildCategory.Close()
		for rowsChildCategory.Next() {
			var childCategory models.Category
			if err := rowsChildCategory.Scan(&childCategory.ID, &childCategory.NameTM, &childCategory.NameRU, &childCategory.ParentCategoryID); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			// childyn child categorysy alynyar
			rowsChildChildCategory, err := db.Query(context.Background(), queryForChildCategory, shopID, childCategory.ID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			defer rowsChildChildCategory.Close()
			for rowsChildChildCategory.Next() {
				var childchildCategory models.Category
				if err := rowsChildChildCategory.Scan(&childchildCategory.ID, &childchildCategory.NameTM, &childchildCategory.NameRU, &childchildCategory.ParentCategoryID); err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}
				childCategory.ChildCategories = append(childCategory.ChildCategories, childchildCategory)
			}
			category.ChildCategories = append(category.ChildCategories, childCategory)
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"categories": categories,
	})
}
