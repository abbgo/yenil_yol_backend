package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	var product models.Product
	db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,price,old_price,brend_id FROM products WHERE id=$1 AND deleted_at IS NULL", productID).Scan(&product.ID, &product.NameTM, &product.NameRU, &product.Price, &product.OldPrice, &product.BrendID)
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
			var image string
			if err := rowsImage.Scan(&image); err != nil {
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

func GetProducts(c *gin.Context) {
	var products []models.Product
	requestQuery := models.ProductQuery{StandartQuery: helpers.StandartQuery{IsDeleted: false}}
	var count uint

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

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	countQuery := "SELECT COUNT(DISTINCT(p.id)) FROM products p INNER JOIN category_products cp ON cp.product_id=p.id WHERE cp.category_id=ANY($1) AND p.deleted_at IS NULL AND cp.deleted_at IS NULL"
	if requestQuery.ShopID != "" {
		rows := strings.Split(countQuery, " WHERE ")
		countQuery = fmt.Sprintf("%v INNER JOIN shop_categories sc ON sc.category_id=cp.category_id WHERE sc.shop_id='%v' AND sc.deleted_at IS NULL AND %v ", rows[0], requestQuery.ShopID, rows[1])
	}
	// database - den product - laryn sany alynyar
	if err = db.QueryRow(context.Background(), countQuery, pq.Array(requestQuery.Categories)).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora product - lary almak ucin query yazylyar
	rowQuery := "SELECT DISTINCT ON (p.id,p.created_at) p.id,p.name_tm,p.name_ru,p.price,p.old_price FROM products p INNER JOIN category_products cp ON cp.product_id=p.id WHERE cp.category_id=ANY($1) AND p.deleted_at IS NULL AND cp.deleted_at IS NULL ORDER BY p.created_at DESC LIMIT $2 OFFSET $3"
	if requestQuery.ShopID != "" {
		rows := strings.Split(rowQuery, " WHERE ")
		rowQuery = fmt.Sprintf("%v INNER JOIN shop_categories sc ON sc.category_id=cp.category_id WHERE sc.shop_id='%v' AND sc.deleted_at IS NULL AND %v ", rows[0], requestQuery.ShopID, rows[1])
	}

	// product - lar alynyar
	rowsProducts, err := db.Query(context.Background(), rowQuery, pq.Array(requestQuery.Categories), requestQuery.Limit, offset)
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
		if err := db.QueryRow(context.Background(), "SELECT pi.image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id WHERE pc.product_id=$1 AND pi.deleted_at IS NULL AND pc.deleted_at IS NULL LIMIT 1", product.ID).Scan(&product.Image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
		"total":    count,
	})
}
