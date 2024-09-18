package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	modelHelpers "github/abbgo/yenil_yol/backend/helpers/models"
	"github/abbgo/yenil_yol/backend/models"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gopkg.in/guregu/null.v4"
)

func CreateCategory(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCategory(category, true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	var parent_category_id interface{}
	if category.ParentCategoryID.String == "" {
		parent_category_id = nil
	} else {
		parent_category_id = category.ParentCategoryID.String
	}

	// eger maglumatlar dogry bolsa onda categories tablisa maglumatlar gosulyar
	_, err = db.Exec(
		context.Background(),
		`INSERT INTO categories (name_tm,name_ru,image,slug_tm,slug_ru,dimension_group_id,parent_category_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		category.NameTM, category.NameRU, category.Image,
		slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"),
		category.DimensionGroupID, parent_category_id,
	)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// category - nyn maglumatlary gosulandan sonra suraty bar bolsa helper_images tablisa category ucin gosulan surat pozulyar
	if category.Image.String != "" {
		if err := DeleteImageFromDB(category.Image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}

func UpdateCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var category models.Category
	if err := c.BindJSON(&category); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCategory(category, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(),
		"UPDATE categories SET name_tm=$1 , name_ru=$2 , image=$3 , slug_tm=$4 , slug_ru=$5, dimension_group_id=$6, parent_category_id=$7 WHERE id=$8",
		category.NameTM, category.NameRU, category.Image, slug.MakeLang(category.NameTM, "en"), slug.MakeLang(category.NameRU, "en"), category.DimensionGroupID, category.ParentCategoryID, category.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if category.Image.String != "" {
		if err := DeleteImageFromDB(category.Image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func GetCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden category id alynyar
	categoryID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var category serializations.GetCategoriesForAdmin
	if err := db.QueryRow(
		context.Background(),
		"SELECT id,name_tm,name_ru,image,dimension_group_id,parent_category_id FROM categories WHERE id = $1 AND deleted_at IS NULL",
		categoryID,
	).
		Scan(&category.ID, &category.NameTM, &category.NameRU, &category.Image, &category.DimensionGroupID, &category.ParentCategoryID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	//kategoriya degisli razmer grupba we razmerler alynyar
	db.QueryRow(context.Background(), `SELECT id,name FROM dimension_groups WHERE id=$1`, category.DimensionGroupID).
		Scan(&category.DimensionGroup.ID, &category.DimensionGroup.Name)
	rowsDimensions, err := db.Query(context.Background(), `SELECT dimension FROM dimensions WHERE dimension_group_id=$1`, category.DimensionGroupID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsDimensions.Close()

	for rowsDimensions.Next() {
		var dimension string
		if err := rowsDimensions.Scan(&dimension); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		category.DimensionGroup.Dimensions = append(category.DimensionGroup.Dimensions, dimension)
	}

	// eger kategoriyanyn parenti bar bolsa db - den parent category alynyar
	if category.ParentCategoryID.String != "" {
		var parentCategory serializations.CategoryForProduct
		if err := db.QueryRow(context.Background(), `SELECT id,name_tm,name_ru FROM categories WHERE id=$1`, category.ParentCategoryID.String).
			Scan(&parentCategory.ID, &parentCategory.NameTM, &parentCategory.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		category.ParentCategory = &parentCategory
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if category.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"category": category,
	})
}

func GetCategoriesWithChild(c *gin.Context) {
	var categories []serializations.GetCategoriesForAdmin
	requestQuery := serializations.CategoryQuery{}
	var searchQuery, search, searchStr, parentCategoryQuery string
	count := 0

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

	if requestQuery.Search != "" {
		incomingsSarch := slug.MakeLang(c.Query("search"), "en")
		search = strings.ReplaceAll(incomingsSarch, "-", " | ")
		searchStr = fmt.Sprintf("%%%s%%", search)
	}

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	orderByQuery := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	if requestQuery.Search != "" {
		searchQuery = fmt.Sprintf(` %s (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s') `, `AND`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	} else {
		parentCategoryQuery = `AND parent_category_id IS NULL`
	}

	// db - den maglumatlaryn sany alynyar
	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL %s %s `, parentCategoryQuery, searchQuery)
	if err := db.QueryRow(context.Background(), queryCount).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// db - den maglumatlar alynyar
	rowQuery := fmt.Sprintf(
		`SELECT id,name_tm,name_ru,image,dimension_group_id FROM categories WHERE deleted_at IS NULL %s %s %s`,
		parentCategoryQuery, searchQuery, orderByQuery,
	)

	// shop - a degisli category - ler alynyar
	rowsCategory, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()

	for rowsCategory.Next() {
		var category serializations.GetCategoriesForAdmin
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &category.NameRU, &category.Image, &category.DimensionGroupID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// kategoriya degisli razmer grupbasy alynyar
		category.DimensionGroup, err = modelHelpers.GetDimensionsByDimensionGroupID(category.DimensionGroupID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// child category alynyar
		queryForChildCategory := "SELECT id,name_tm,name_ru,parent_category_id,image,dimension_group_id FROM categories \n\t\tWHERE deleted_at IS NULL AND parent_category_id=$1"

		rowsChildCategory, err := db.Query(context.Background(), queryForChildCategory, category.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsChildCategory.Close()

		for rowsChildCategory.Next() {
			var childCategory serializations.GetCategoriesForAdmin
			if err := rowsChildCategory.Scan(
				&childCategory.ID, &childCategory.NameTM, &childCategory.NameRU,
				&childCategory.ParentCategoryID, &childCategory.Image, &childCategory.DimensionGroupID,
			); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			// child kategoriya degisli razmer grupbasy alynyar
			childCategory.DimensionGroup, err = modelHelpers.GetDimensionsByDimensionGroupID(childCategory.DimensionGroupID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			// childyn child categorysy alynyar
			rowsChildChildCategory, err := db.Query(context.Background(), queryForChildCategory, childCategory.ID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			defer rowsChildChildCategory.Close()
			for rowsChildChildCategory.Next() {
				var childchildCategory serializations.GetCategoriesForAdmin
				if err := rowsChildChildCategory.Scan(
					&childchildCategory.ID, &childchildCategory.NameTM,
					&childchildCategory.NameRU, &childchildCategory.ParentCategoryID, &childchildCategory.Image,
					&childchildCategory.DimensionGroupID,
				); err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}

				// child child - yn kategoriya degisli razmer grupbasy alynyar
				childchildCategory.DimensionGroup, err = modelHelpers.GetDimensionsByDimensionGroupID(childchildCategory.DimensionGroupID)
				if err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}

				childCategory.ChildCategories = append(childCategory.ChildCategories, childchildCategory)
			}
			category.ChildCategories = append(category.ChildCategories, childCategory)
		}
		categories = append(categories, category)
	}

	pageCount := count / requestQuery.Limit
	if count%requestQuery.Limit != 0 {
		pageCount = count/requestQuery.Limit + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"categories": categories,
		"count":      count,
		"page_count": pageCount,
	})
}

func GetDeletedCategories(c *gin.Context) {
	var categories []serializations.GetCategoriesForAdmin
	requestQuery := serializations.CategoryQuery{}
	var searchQuery, search, searchStr string
	count := 0

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

	if requestQuery.Search != "" {
		incomingsSarch := slug.MakeLang(c.Query("search"), "en")
		search = strings.ReplaceAll(incomingsSarch, "-", " | ")
		searchStr = fmt.Sprintf("%%%s%%", search)
	}

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	orderByQuery := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	if requestQuery.Search != "" {
		searchQuery = fmt.Sprintf(` %s (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s') `, `AND`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}

	// db - den maglumatlaryn sany alynyar
	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM categories WHERE deleted_at IS NOT NULL %s `, searchQuery)
	if err := db.QueryRow(context.Background(), queryCount).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// db - den maglumatlar alynyar
	rowQuery := fmt.Sprintf(
		`SELECT id,name_tm,name_ru,image,dimension_group_id FROM categories WHERE deleted_at IS NOT NULL %s %s`,
		searchQuery, orderByQuery,
	)

	// shop - a degisli category - ler alynyar
	rowsCategory, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()

	for rowsCategory.Next() {
		var category serializations.GetCategoriesForAdmin
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &category.NameRU, &category.Image, &category.DimensionGroupID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// kategoriya degisli razmer grupbasy alynyar
		category.DimensionGroup, err = modelHelpers.GetDimensionsByDimensionGroupID(category.DimensionGroupID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		categories = append(categories, category)
	}

	pageCount := count / requestQuery.Limit
	if count%requestQuery.Limit != 0 {
		pageCount = count/requestQuery.Limit + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"categories": categories,
		"count":      count,
		"page_count": pageCount,
	})
}

func GetCategories(c *gin.Context) {
	var categories []serializations.CategoryForProduct
	requestQuery := serializations.CategoryQuery{}
	var searchQuery, search, searchStr string
	count := 0

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

	if requestQuery.Search != "" {
		incomingsSarch := slug.MakeLang(c.Query("search"), "en")
		search = strings.ReplaceAll(incomingsSarch, "-", " | ")
		searchStr = fmt.Sprintf("%%%s%%", search)
	}

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	orderByQuery := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	if requestQuery.Search != "" {
		searchQuery = fmt.Sprintf(` %s (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s') `, `AND`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}

	// db - den maglumatlaryn sany alynyar
	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL %s `, searchQuery)
	if err := db.QueryRow(context.Background(), queryCount).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// db - den maglumatlar alynyar
	rowQuery := fmt.Sprintf(
		`SELECT id,name_tm,name_ru FROM categories WHERE deleted_at IS NULL %s %s`,
		searchQuery, orderByQuery,
	)

	// shop - a degisli category - ler alynyar
	rowsCategory, err := db.Query(context.Background(), rowQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()

	for rowsCategory.Next() {
		var category serializations.CategoryForProduct
		if err := rowsCategory.Scan(&category.ID, &category.NameTM, &category.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		categories = append(categories, category)
	}

	pageCount := count / requestQuery.Limit
	if count%requestQuery.Limit != 0 {
		pageCount = count/requestQuery.Limit + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"categories": categories,
		"count":      count,
		"page_count": pageCount,
	})
}

func DeleteCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den category id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("categories", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa categories we sol tablisa bilen baglanysykly tablisalaryn deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_category($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den category id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("categories", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa categories we sol tablisa bilen baglanysykly tablisalaryn deleted_at - ine NULL goyulyar
	_, err = db.Exec(context.Background(), "CALL restore_category($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyCategoryByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den category id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	var image sql.NullString
	db.QueryRow(context.Background(), "SELECT id,image FROM categories WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id, &image)

	// eger database - de gelen id degisli category yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger category - nyn suraty bar bolsa onda ol local papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := os.Remove(helpers.ServerPath + "assets/" + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// category - nyn suraty pozulandan son category we onun bilen baglanysykly maglumatlar pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func CheckForDelete(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden category id alynyar
	categoryID := c.Param("id")

	// Ilki bilen kategoriyanyn barlgygy barlanyar
	if err := helpers.ValidateRecordByID("categories", categoryID, "NULL", db); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// Kategoriya degisli child categoriya barmy sol barlanyar
	var countOfChildCategories uint8
	if err := db.QueryRow(
		context.Background(), `SELECT COUNT(id) FROM categories WHERE parent_category_id=$1 AND deleted_at IS NULL`, categoryID,
	).
		Scan(&countOfChildCategories); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// Kategoriya degisli haryt barmy sol barlanyar
	var countOfProducts uint8
	if err := db.QueryRow(
		context.Background(), `SELECT COUNT(id) FROM category_products WHERE category_id=$1 AND deleted_at IS NULL`, categoryID,
	).
		Scan(&countOfProducts); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"for_deletion": countOfProducts == 0 && countOfChildCategories == 0,
	})
}

func GetParentCategory(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden category id alynyar
	categoryID := c.Param("id")

	// Ilki bilen kategoriyanyn barlgygy barlanyar
	var id string
	var parentCategoryID null.String
	db.QueryRow(context.Background(), `SELECT id,parent_category_id FROM categories WHERE id=$1`, categoryID).
		Scan(&id, &parentCategoryID)
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// Kategoriya degisli pozulan parent category barmy yada yok sol barlanyar
	var parentCategory serializations.CategoryForProduct
	if parentCategoryID.String != "" {

		if err := db.QueryRow(
			context.Background(), `SELECT id,name_tm,name_ru FROM categories WHERE id=$1 AND deleted_at IS NOT NULL`,
			parentCategoryID.String,
		).
			Scan(&parentCategory.ID, &parentCategory.NameTM, &parentCategory.NameRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          true,
		"parent_category": parentCategory,
	})
}
