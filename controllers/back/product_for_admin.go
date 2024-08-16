package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
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
	}

	// request - den gelen deleted statusa gora pozulan ya-da pozulmadyk maglumatlar alynmaly
	if requestQuery.IsDeleted {
		isDeleted = "NOT NULL"
	}

	// request query - den status - a gora product - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(`SELECT id,name_tm,name_ru,price,old_price,brend_id,shop_id,is_visible FROM products p WHERE deleted_at IS %v`, isDeleted)
	orderQuery := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	if requestQuery.Search != "" {
		searchQuery = fmt.Sprintf(` %s (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s') `, `AND`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}

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

		// haryda degisli suratlar alynyar
		// if err := db.QueryRow(context.Background(),
		// 	`SELECT image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id WHERE pc.product_id=$1 AND pc.order_number=1 AND pi.order_number=1`,
		// 	product.ID).Scan(&product.Image); err != nil {
		// 	helpers.HandleError(c, 400, err.Error())
		// 	return
		// }
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
	})
}
