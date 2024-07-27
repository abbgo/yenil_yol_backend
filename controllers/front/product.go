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

func GetProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	productID := c.Param("id")

	// request - dan gelen id boyunca haryt alynyar
	var product serializations.GetProductsForFront
	db.QueryRow(context.Background(),
		`SELECT DISTINCT ON (p.id) p.id,p.name_tm,p.name_ru,p.price,p.old_price,p.brend_id,s.id,s.name_tm,s.name_ru FROM products p 
		INNER JOIN shops s ON s.id=p.shop_id
		WHERE p.id=$1 AND p.deleted_at IS NULL AND p.is_visible=true`,
		productID).Scan(
		&product.ID,
		&product.NameTM,
		&product.NameRU,
		&product.Price,
		&product.OldPrice,
		&product.BrendID,
		&product.Shop.ID,
		&product.Shop.NameTM,
		&product.Shop.NameRU)
	if product.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// harydyn brendi alynyar
	db.QueryRow(context.Background(), "SELECT name,image FROM brends WHERE id=$1 AND deleted_at IS NULL", product.BrendID).Scan(&product.Brend.Name, &product.Brend.Image)

	// sonra harydyn renkleri we ona degisli suratlar alynyar
	rowsColor, err := db.Query(context.Background(), "SELECT id FROM product_colors WHERE product_id=$1", product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsColor.Close()
	for rowsColor.Next() {
		var productColor models.ProductColor
		if err := rowsColor.Scan(&productColor.ID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// sonra renke degisli razmerler alynyar
		rowsDimension, err := db.Query(context.Background(), "SELECT d.dimension FROM dimensions d INNER JOIN product_dimensions pd ON pd.dimension_id=d.id WHERE d.deleted_at IS NULL AND pd.deleted_at IS NULL AND pd.product_color_id=$1", productColor.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsDimension.Close()
		for rowsDimension.Next() {
			var dimension string
			if err := rowsDimension.Scan(&dimension); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			productColor.Dimensions = append(productColor.Dimensions, dimension)
		}

		// sonra renke degisli suratlar alynyar
		rowsImage, err := db.Query(context.Background(), "SELECT image FROM product_images WHERE product_color_id=$1 AND deleted_at IS NULL", productColor.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsColor.Close()
		for rowsImage.Next() {
			var image models.ProductImage
			if err := rowsImage.Scan(&image.Image); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			productColor.Images = append(productColor.Images, image)
		}

		product.ProductColors = append(product.ProductColors, productColor)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"product": product,
	})
}

func GetSimilarProductsByProductID(c *gin.Context) {
	var products []models.Product
	requestQuery := serializations.ProductQuery{StandartQuery: helpers.StandartQuery{IsDeleted: false}}

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

	// request query - den maglumatlara gora product - lary almak ucin query yazylyar
	rowQuery := `SELECT DISTINCT ON (p.id,p.created_at) p.id,p.name_tm,p.name_ru,p.price,p.old_price FROM products p 
				INNER JOIN category_products cp ON cp.product_id=p.id 
				INNER JOIN categories c ON c.id=cp.category_id WHERE p.id!=$1
				AND c.id = (SELECT ca.id FROM categories ca INNER JOIN category_products cap ON cap.category_id=ca.id WHERE 
				ca.parent_category_id IS NOT NULL AND cap.product_id=$1 ORDER BY ca.created_at DESC LIMIT 1) 
				AND c.deleted_at IS NULL 
				AND p.deleted_at IS NULL 
				AND cp.deleted_at IS NULL 
				AND p.is_visible=true
				ORDER BY p.created_at DESC LIMIT $2
				`

	// product - lar alynyar
	rowsProducts, err := db.Query(context.Background(), rowQuery, requestQuery.ProductID, requestQuery.Limit)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	for rowsProducts.Next() {
		var product models.Product
		if err := rowsProducts.Scan(&product.ID, &product.NameTM, &product.NameRU, &product.Price, &product.OldPrice); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// haryda degisli yekeje surat alyas
		if err := db.QueryRow(context.Background(),
			`SELECT pi.image FROM product_images pi 
			INNER JOIN product_colors pc ON pc.id=pi.product_color_id 
			WHERE pc.product_id=$1 AND pi.deleted_at IS NULL 
			AND pc.deleted_at IS NULL LIMIT 1`,
			product.ID).Scan(&product.Image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
	})

}

func GetProductsByIDs(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden product id - ler alynyar
	productIDs := c.QueryArray("ids")

	// database - den request parametr - den gelen id - ler boyunca maglumat cekilyar
	var products []models.Product
	rows, err := db.Query(context.Background(),
		`
			SELECT id,name_tm,name_ru,price,old_price FROM products 
			WHERE id = ANY($1) AND deleted_at IS NULL AND is_visible=true
		`,
		pq.Array(productIDs))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.NameTM, &product.NameRU, &product.Price, &product.OldPrice); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// haryda degisli yekeje surat alyas
		db.QueryRow(
			context.Background(), `
								SELECT pi.resized_image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id 
								WHERE pc.product_id=$1 AND pc.order_number=1 AND pi.order_number=1 AND pi.deleted_at IS NULL AND pc.deleted_at IS NULL
							`,
			product.ID).Scan(&product.Image)

		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
	})
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	requestQuery := serializations.ProductQuery{StandartQuery: helpers.StandartQuery{IsDeleted: false}}
	var shopWhereQuery, categoryJoinQuery, categoryQuery, searchQuery, search, searchStr string
	isVisibleQuery := ` WHERE p.is_visible=true `

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

	// request query - den status - a gora product - lary almak ucin query yazylyar
	defaultQuery := `SELECT DISTINCT ON (p.id,p.created_at) p.id,p.name_tm,p.name_ru,p.price,p.old_price FROM products p`

	if requestQuery.ShopID != "" {
		shopWhereQuery = fmt.Sprintf(` AND p.shop_id='%s' `, requestQuery.ShopID)
	}

	if requestQuery.Search != "" {
		isVisibleQuery = ""
		searchQuery = fmt.Sprintf(` %s (to_tsvector(p.slug_%s) @@ to_tsquery('%s') OR p.slug_%s LIKE '%s') `, `WHERE p.is_visible=true AND `, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}

	if len(requestQuery.Categories) != 0 {
		categoryJoinQuery = ` INNER JOIN category_products cp ON cp.product_id=p.id `
		categoryQuery = ` AND cp.category_id=ANY($1) AND p.deleted_at IS NULL AND p.is_visible=true AND cp.deleted_at IS NULL `
		if requestQuery.Search != "" {
			searchQuery = fmt.Sprintf(` %s (to_tsvector(p.slug_%s) @@ to_tsquery('%s') OR p.slug_%s LIKE '%s') `, `AND`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
		}
	}

	// product - lar alynyar
	var rowsProducts pgx.Rows
	if len(requestQuery.Categories) != 0 {
		rowsProducts, err = db.Query(context.Background(), defaultQuery+categoryJoinQuery+isVisibleQuery+categoryQuery+shopWhereQuery+searchQuery+` ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`, pq.Array(requestQuery.Categories), requestQuery.Limit, offset)
	} else {
		rowsProducts, err = db.Query(context.Background(), defaultQuery+isVisibleQuery+searchQuery+shopWhereQuery+` ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`, requestQuery.Limit, offset)
	}
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	for rowsProducts.Next() {
		var product models.Product
		if err := rowsProducts.Scan(&product.ID, &product.NameTM, &product.NameRU, &product.Price, &product.OldPrice); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// haryda degisli yekeje surat alyas
		err := db.QueryRow(context.Background(), `SELECT DISTINCT ON (pi.id) pi.resized_image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id WHERE pc.product_id=$1 AND pc.order_number=1 AND pi.order_number=1 AND pi.deleted_at IS NULL AND pc.deleted_at IS NULL`, product.ID).Scan(&product.Image)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
	})
}
