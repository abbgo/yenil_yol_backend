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
		queryForChildCategory := `SELECT DISTINCT ON (c.id) c.id,c.name_tm,c.name_ru FROM categories c 
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
			if err := rowsChildCategory.Scan(&childCategory.ID, &childCategory.NameTM, &childCategory.NameRU); err != nil {
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
				if err := rowsChildChildCategory.Scan(&childchildCategory.ID, &childchildCategory.NameTM, &childchildCategory.NameRU); err != nil {
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

// func GetCategoriesByCategoryID(c *gin.Context) {
// 	var requestQuery = models.CategoryQuery{}
// 	var categories []models.Category

// 	// request query - den maglumatlar bind edilyar
// 	if err := c.Bind(&requestQuery); err != nil {
// 		helpers.HandleError(c, 400, err.Error())
// 		return
// 	}
// 	// request query - den maglumatlar validate edilyar
// 	if err := helpers.ValidateStructData(&requestQuery); err != nil {
// 		helpers.HandleError(c, 400, err.Error())
// 		return
// 	}

// 	// initialize database connection
// 	db, err := config.ConnDB()
// 	if err != nil {
// 		helpers.HandleError(c, 400, err.Error())
// 		return
// 	}
// 	defer db.Close()

// 	// shop - a degisli category - ler alynyar
// 	queryCategoriesByShopID := fmt.Sprintf(`
// 		SELECT c.id,c.name_tm,c.name_ru FROM categories c
// 		INNER JOIN shop_categories sc ON sc.category_id=c.id
// 		WHERE sc.shop_id=$1 AND c.parent_category_id IS NULL AND c.deleted_at IS NULL AND sc.deleted_at IS NULL
// 	`)

// 	rowsCategoriesByShopID, err := db.Query(context.Background(),
// 		`SELECT c.id,c.name_tm,c.name_ru FROM categories c
// 		INNER JOIN shop_categories sc ON sc.category_id=c.id
// 		WHERE sc.shop_id=$1 AND c.parent_category_id IS NULL AND c.deleted_at IS NULL AND sc.deleted_at IS NULL`,
// 		requestQuery.ShopID)
// 	if err != nil {
// 		helpers.HandleError(c, 400, err.Error())
// 		return
// 	}
// 	defer rowsCategoriesByShopID.Close()
// 	for rowsCategoriesByShopID.Next() {
// 		var parentCategory models.Category
// 		if err := rowsCategoriesByShopID.Scan(&parentCategory.ID, &parentCategory.NameTM, &parentCategory.NameRU); err != nil {
// 			helpers.HandleError(c, 400, err.Error())
// 			return
// 		}
// 		// categories = append(categories, category)

// 		// child category alynyar
// 		rowsChildCategory, err := db.Query(context.Background(),
// 			`SELECT c.id,c.name_tm,c.name_ru FROM categories c
// 		INNER JOIN shop_categories sc ON sc.category_id=c.id
// 		WHERE sc.shop_id=$1 AND c.parent_category_id=$2 AND c.deleted_at IS NULL AND sc.deleted_at IS NULL`,
// 			requestQuery.ShopID, requestQuery.CategoryID)
// 		if err != nil {
// 			helpers.HandleError(c, 400, err.Error())
// 			return
// 		}
// 		defer rowsChildCategory.Close()
// 		for rowsChildCategory.Next() {
// 			var childCategory models.Category
// 			if err := rowsChildCategory.Scan(&childCategory.ID, &childCategory.NameTM, &childCategory.NameRU); err != nil {
// 				helpers.HandleError(c, 400, err.Error())
// 				return
// 			}
// 			category.ChildCategories = append(category.ChildCategories, childCategory)
// 		}

// 	}

// 	// var category models.Category
// 	// db.QueryRow(context.Background(),
// 	// 	`SELECT id,name_tm,name_ru FROM categories WHERE id=$1 AND deleted_at IS NULL`,
// 	// 	requestQuery.CategoryID).Scan(&category.ID, &category.NameTM, &category.NameRU)
// 	// if category.ID == "" {
// 	// 	helpers.HandleError(c, 404, "record not found")
// 	// 	return
// 	// }

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":   true,
// 		"category": category,
// 	})
// }
