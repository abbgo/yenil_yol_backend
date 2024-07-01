package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

func CreateProduct(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	productCode, err := models.ValidateProduct(product, true)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda products tablisa maglumatlar gosulyar we yzyna id return edilyar
	if err := db.QueryRow(context.Background(), "INSERT INTO products (name_tm,name_ru,price,old_price,code,slug_tm,slug_ru,brend_id,is_visible) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id", product.NameTM, product.NameRU, product.Price, product.OldPrice, productCode, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.BrendID, product.IsVisible).Scan(&product.ID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// haryt gosulandan son category_products tablisa haryda degisli category - ler gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO category_products (product_id,category_id) VALUES ($1,unnest($2::uuid[]))", product.ID, pq.Array(product.Categories))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// bu yerde harydyn renkleri we sol renklere degisli suratlar we razmerler gosulyar
	for _, v := range product.ProductColors {
		var productColorID string
		if err := db.QueryRow(context.Background(), "INSERT INTO product_colors (name,product_id) VALUES ($1,$2) RETURNING id", v.Name, product.ID).Scan(&productColorID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// renk gosulandan sonra sol renke degisli razmerler gosulyar
		_, err := db.Exec(context.Background(), "INSERT INTO product_dimensions (product_color_id,dimension_id) VALUES ($1,unnest($2::uuid[]))", productColorID, pq.Array(v.Dimensions))
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// bu yerde renke degisli suratlar gosulyar
		_, err = db.Exec(context.Background(), "INSERT INTO product_images (product_color_id,image) VALUES ($1,unnest($2::varchar[]))", productColorID, pq.Array(v.Images))
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}

func UpdateProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	productCode, err := models.ValidateProduct(product, false)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE products SET name_tm=$1 , name_ru=$2 , price=$3 , old_price=$4 , code=$5 , slug_tm=$6 , slug_ru=$7 , brend_id=$8 , is_visible=$9 WHERE id=$10", product.NameTM, product.NameRU, product.Price, product.OldPrice, productCode, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.BrendID, product.IsVisible, product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// harydyn maglumatlary uytgedilenson degisli edilen onki category - ler pozulyp tazeleri gosulyar
	_, err = db.Exec(context.Background(), "DELETE FROM category_products WHERE product_id=$1", product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	_, err = db.Exec(context.Background(), "INSERT INTO category_products (product_id,category_id) VALUES ($1,unnest($2::uuid[]))", product.ID, pq.Array(product.Categories))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// bu yerde haryda degisli renkler , razmerler we suratlar pozulup olara derek tazesi yazylyar
	_, err = db.Exec(context.Background(), "DELETE FROM product_dimensions pd USING product_colors pc WHERE pc.id=pd.product_color_id AND pc.product_id=$1", product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	_, err = db.Exec(context.Background(), "DELETE FROM product_images pi USING product_colors pc WHERE pc.id=pi.product_color_id AND pc.product_id=$1", product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	_, err = db.Exec(context.Background(), "DELETE FROM product_colors WHERE product_id=$1", product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	for _, v := range product.ProductColors {
		var productColorID string
		if err := db.QueryRow(context.Background(), "INSERT INTO product_colors (name,product_id) VALUES ($1,$2) RETURNING id", v.Name, product.ID).Scan(&productColorID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// renk gosulandan sonra sol renke degisli razmerler gosulyar
		_, err := db.Exec(context.Background(), "INSERT INTO product_dimensions (product_color_id,dimension_id) VALUES ($1,unnest($2::uuid[]))", productColorID, pq.Array(v.Dimensions))
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// bu yerde renke degisli suratlar gosulyar
		_, err = db.Exec(context.Background(), "INSERT INTO product_images (product_color_id,image) VALUES ($1,unnest($2::varchar[]))", productColorID, pq.Array(v.Images))
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func GetProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden product id alynyar
	productID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var product models.Product
	if err := db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,price,old_price,code,brend_id,is_visible FROM products WHERE id = $1 AND deleted_at IS NULL", productID).Scan(
		&product.ID,
		&product.NameTM,
		&product.NameRU,
		&product.Price,
		&product.OldPrice,
		&product.Code,
		&product.BrendID,
		&product.IsVisible,
	); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// haryda degisli category - lar alynyar
	rowsCategory, err := db.Query(context.Background(), "SELECT category_id FROM category_products WHERE product_id=$1 AND deleted_at IS NULL", productID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCategory.Close()

	for rowsCategory.Next() {
		var categoryID string
		if err := rowsCategory.Scan(&categoryID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		product.Categories = append(product.Categories, categoryID)
	}

	// haryda degisli renkler , razmerler we suratlar alynyar
	rowsColor, err := db.Query(context.Background(), "SELECT id,name FROM product_colors WHERE product_id=$1 AND deleted_at IS NULL", productID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsColor.Close()

	for rowsColor.Next() {
		var productColor models.ProductColor
		if err := rowsColor.Scan(&productColor.ID, &productColor.Name); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// renk alynandan son sol renke degisli razmerler alynyar
		rowsDimension, err := db.Query(context.Background(), "SELECT dimension_id FROM product_dimensions WHERE product_color_id=$1 AND deleted_at IS NULL", productColor.ID)
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

		// renk alynandan son sol renke degisli suratlar alynyar
		rowsImage, err := db.Query(context.Background(), "SELECT image FROM product_images WHERE product_color_id=$1 AND deleted_at IS NULL", productColor.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsImage.Close()

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

	// eger databse sol maglumat yok bolsa error return edilyar
	if product.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"product": product,
	})
}

func GetProducts(c *gin.Context) {
	var requestQuery models.ProductQuery
	var count uint
	var products []models.Product
	isDeleted := "NULL"

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

	// request - den gelen deleted statusa gora pozulan ya-da pozulmadyk maglumatlar alynmaly
	if requestQuery.IsDeleted {
		isDeleted = "NOT NULL"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(p.id) FROM products p WHERE p.deleted_at IS %v", isDeleted)
	if requestQuery.ShopID != "" {
		rows := strings.Split(countQuery, " WHERE ")
		countQuery = fmt.Sprintf("%v INNER JOIN category_products cp ON cp.product_id=p.id INNER JOIN shop_categories sc ON sc.category_id=cp.category_id WHERE sc.shop_id='%v' AND sc.deleted_at IS %v AND cp.deleted_at IS %v AND %v ", rows[0], requestQuery.ShopID, isDeleted, isDeleted, rows[1])
	}

	// database - den product - laryn sany alynyar
	if err = db.QueryRow(context.Background(), countQuery).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora product - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf("SELECT p.id,p.name_tm,p.name_ru,p.price,p.old_price,p.code,p.brend_id,p.is_visible FROM products p WHERE p.deleted_at IS %v ORDER BY p.created_at DESC LIMIT $1 OFFSET $2", isDeleted)
	if requestQuery.ShopID != "" {
		rows := strings.Split(rowQuery, " WHERE ")
		rowQuery = fmt.Sprintf("%v INNER JOIN category_products cp ON cp.product_id=p.id INNER JOIN shop_categories sc ON sc.category_id=cp.category_id WHERE sc.shop_id='%v' AND sc.deleted_at IS %v AND cp.deleted_at IS %v AND %v ", rows[0], requestQuery.ShopID, isDeleted, isDeleted, rows[1])
	}

	// database - den brend - lar alynyar
	rowsBrend, err := db.Query(context.Background(), rowQuery, requestQuery.Limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsBrend.Close()

	for rowsBrend.Next() {
		var product models.Product
		if err := rowsBrend.Scan(
			&product.ID,
			&product.NameTM,
			&product.NameRU,
			&product.Price,
			&product.OldPrice,
			&product.Code,
			&product.BrendID,
			&product.IsVisible,
		); err != nil {
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

func DeleteProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den product id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("products", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	_, err = db.Exec(context.Background(), "CALL delete_product($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den product id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("products", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa brend restore edilyar
	_, err = db.Exec(context.Background(), "CALL restore_product($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den product id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("products", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	rows, err := db.Query(context.Background(), "SELECT pi.image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id WHERE pc.product_id=$1 GROUP BY pi.image", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var image string
		if err := rows.Scan(&image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := os.Remove(helpers.ServerPath + image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := os.Remove(helpers.ServerPath + "assets/" + image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// brend - in suraty pozulandan sonra database - den brend pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM products WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}
