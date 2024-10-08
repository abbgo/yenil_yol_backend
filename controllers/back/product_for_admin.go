package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func GetAdminProducts(c *gin.Context) {
	var requestQuery serializations.ProductQuery
	var products []serializations.GetProductsForAdminProduct
	isDeleted := "NULL"
	var searchQuery, search, searchStr string
	count := 0

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

	if requestQuery.Search != "" {
		incomingsSarch := slug.MakeLang(c.Query("search"), "en")
		search = strings.ReplaceAll(incomingsSarch, "-", " | ")
		searchStr = fmt.Sprintf("%%%s%%", search)
		searchQuery = fmt.Sprintf(` %s (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s') `, `AND`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}

	// request - den gelen deleted statusa gora pozulan ya-da pozulmadyk maglumatlar alynmaly
	if requestQuery.IsDeleted {
		isDeleted = "NOT NULL"
	}

	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM products WHERE deleted_at IS %v `, isDeleted)
	if len(requestQuery.CratedStatuses) != 0 {
		if err := db.QueryRow(
			context.Background(), queryCount+searchQuery+" AND created_status=ANY($1) ",
			pq.Array(requestQuery.CratedStatuses),
		).Scan(&count); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	} else {
		if err := db.QueryRow(context.Background(), queryCount+searchQuery).Scan(&count); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// request query - den status - a gora product - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(`SELECT id,name_tm,name_ru,price,old_price,brend_id,shop_id,is_visible,genders FROM products WHERE deleted_at IS %v`, isDeleted)
	orderQuery := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	var rowsProducts pgx.Rows
	if len(requestQuery.CratedStatuses) != 0 {
		rowsProducts, err = db.Query(context.Background(), rowQuery+searchQuery+" AND created_status=ANY($1) "+orderQuery, pq.Array(requestQuery.CratedStatuses))
	} else {
		rowsProducts, err = db.Query(context.Background(), rowQuery+searchQuery+orderQuery)
	}
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsProducts.Close()

	for rowsProducts.Next() {
		var product serializations.GetProductsForAdminProduct
		if err := rowsProducts.Scan(
			&product.ID,
			&product.NameTM,
			&product.NameRU,
			&product.Price,
			&product.OldPrice,
			&product.BrendID,
			&product.ShopID,
			&product.IsVisible,
			&product.Genders,
		); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// eger harydyn brendi bar bolsa haryda degisli brend alynyar
		if product.BrendID.String != "" {
			db.QueryRow(context.Background(), `SELECT id,name,image FROM brends WHERE id=$1`, product.BrendID.String).
				Scan(&product.Brend.ID, &product.Brend.Name, &product.Brend.Image)
		}

		// harydyn dukany alynyar
		db.QueryRow(context.Background(), `SELECT id,name_tm,name_ru FROM shops WHERE id=$1`, product.ShopID).
			Scan(&product.Shop.ID, &product.Shop.NameTM, &product.Shop.NameRU)

		// Harydyn kategoriyalary alynyar
		rowsCategories, err := db.Query(
			context.Background(),
			`SELECT DISTINCT ON (c.id) c.id,c.name_tm,c.name_ru FROM categories c 
			INNER JOIN category_products cp ON cp.category_id=c.id 
			WHERE cp.product_id=$1`, product.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsCategories.Close()

		for rowsCategories.Next() {
			var category serializations.CategoryForProduct
			if err := rowsCategories.Scan(&category.ID, &category.NameTM, &category.NameRU); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			product.Categories = append(product.Categories, category)
		}

		// harydyn renkleri alynyar
		rowsColors, err := db.Query(context.Background(), `SELECT id,name FROM product_colors WHERE product_id=$1`, product.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsColors.Close()

		for rowsColors.Next() {
			var color serializations.ProductColorForAdmin
			if err := rowsColors.Scan(&color.ID, &color.Name); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			// harydyn bu renkine degisli razmerler alynyar
			rowsDimensions, err := db.Query(
				context.Background(),
				`SELECT DISTINCT ON (d.dimension) d.dimension FROM dimensions d 
				INNER JOIN product_dimensions pd ON pd.dimension_id=d.id 
				WHERE pd.product_color_id=$1`, color.ID)
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
				color.Dimensions = append(color.Dimensions, dimension)
			}

			// harydyn renkine degisli suratlar alynyar
			rowsImages, err := db.Query(context.Background(), `SELECT image FROM product_images WHERE product_color_id=$1`, color.ID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			defer rowsImages.Close()

			for rowsImages.Next() {
				var image string
				if err := rowsImages.Scan(&image); err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}
				color.Images = append(color.Images, image)
			}

			product.ProductColors = append(product.ProductColors, color)
		}

		products = append(products, product)
	}

	pageCount := count / requestQuery.Limit
	if count%requestQuery.Limit != 0 {
		pageCount = count/requestQuery.Limit + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"products":   products,
		"count":      count,
		"page_count": pageCount,
	})
}

func UpdateProductCreatedStatus(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var product models.UpdateCreatedStatusShop
	if err := c.BindJSON(&product); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateUpdateProductCreatedStatus(product); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	var rejectedReason interface{}
	if product.RejectedReason != "" {
		rejectedReason = product.RejectedReason
	} else {
		rejectedReason = nil
	}

	// maglumatlar barlananda son product - yn created status - y update edilyar
	_, err = db.Exec(context.Background(), `UPDATE products SET created_status=$1 , rejected_reason=$2 WHERE id=$3`, product.CreatedStatus, rejectedReason, product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}
