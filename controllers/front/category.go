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

func GetCategoriesByCategoryID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	shopID := c.Param("shop_id")
	// request parametrden category id alynyar
	categoryID := c.Param("category_id")

	var category models.Category
	db.QueryRow(context.Background(),
		`SELECT id,name_tm,name_ru FROM categories WHERE id=$1 AND deleted_at IS NULL`,
		categoryID).Scan(&category.ID, &category.NameTM, &category.NameRU)
	if category.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	rowsChildCategory, err := db.Query(context.Background(),
		`SELECT c.id,c.name_tm,c.name_ru FROM categories c 
		INNER JOIN shop_categories sc ON sc.category_id=c.id 
		WHERE sc.shop_id=$1 AND c.parent_category_id=$2 AND c.deleted_at IS NULL AND sc.deleted_at IS NULL`,
		shopID, categoryID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsChildCategory.Close()
	for rowsChildCategory.Next() {
		var childCategory models.Category
		if err := rowsChildCategory.Scan(&childCategory.ID, &childCategory.NameTM, &childCategory.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		category.ChildCategories = append(category.ChildCategories, childCategory)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"category": category,
	})
}
