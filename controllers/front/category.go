package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []serializations.GetCategories
	var requestQuery serializations.CategoryQuery

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

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	rowQuery := `SELECT id,name_tm,name_ru FROM categories WHERE deleted_at IS NULL AND parent_category_id IS NULL`
	if requestQuery.ShopID != "" {
		rowQuery = fmt.Sprintf(`SELECT DISTINCT ON (c.id) c.id,c.name_tm,c.name_ru FROM categories c
		INNER JOIN category_products cp ON cp.category_id=c.id
		INNER JOIN products p ON p.id=cp.product_id
		WHERE p.shop_id='%s' AND c.parent_category_id IS NULL 
		AND c.deleted_at IS NULL 
		AND cp.deleted_at IS NULL 
		AND p.deleted_at IS NULL`, requestQuery.ShopID)
	}

	// shop - a degisli category - ler alynyar
	rowsCategory, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()
	for rowsCategory.Next() {
		var category serializations.GetCategories
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &category.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// child category alynyar
		queryForChildCategory := fmt.Sprintf(`SELECT id,name_tm,name_ru,parent_category_id FROM categories 
		WHERE deleted_at IS NULL AND parent_category_id='%s'`, category.ID)

		if requestQuery.ShopID != "" {
			queryForChildCategory = fmt.Sprintf(`SELECT DISTINCT ON (c.id) c.id,c.name_tm,c.name_ru,c.parent_category_id FROM categories c 
		INNER JOIN category_products cp ON cp.category_id=c.id
		INNER JOIN products p ON p.id=cp.product_id
		WHERE p.shop_id='%s' AND c.parent_category_id='%s' 
		AND c.deleted_at IS NULL 
		AND cp.deleted_at IS NULL 
		AND p.deleted_at IS NULL`, requestQuery.ShopID, category.ID)
		}

		rowsChildCategory, err := db.Query(context.Background(), queryForChildCategory)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsChildCategory.Close()
		for rowsChildCategory.Next() {
			var childCategory serializations.GetCategories
			if err := rowsChildCategory.Scan(&childCategory.ID, &childCategory.NameTM, &childCategory.NameRU, &childCategory.ParentCategoryID); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			// childyn child categorysy alynyar
			rowsChildChildCategory, err := db.Query(context.Background(), queryForChildCategory)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			defer rowsChildChildCategory.Close()
			for rowsChildChildCategory.Next() {
				var childchildCategory serializations.GetCategories
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
