package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strconv"

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
	if err := db.QueryRow(context.Background(), "INSERT INTO products (name_tm,name_ru,price,old_price,code,slug_tm,slug_ru,brend_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id", product.NameTM, product.NameRU, product.Price, product.OldPrice, productCode, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.BrendID).Scan(&product.ID); err != nil {
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
	_, err = db.Exec(context.Background(), "UPDATE products SET name_tm=$1 , name_ru=$2 , price=$3 , old_price=$4 , code=$5 , slug_tm=$6 , slug_ru=$7 , brend_id=$8 WHERE id=$9", product.NameTM, product.NameRU, product.Price, product.OldPrice, productCode, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.BrendID, product.ID)
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
	if err := db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,price,old_price,code,brend_id FROM products WHERE id = $1 AND deleted_at IS NULL", productID).Scan(
		&product.ID,
		&product.NameTM,
		&product.NameRU,
		&product.Price,
		&product.OldPrice,
		&product.Code,
		&product.BrendID,
	); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
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
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den limit alynyar
	limitStr := c.Query("limit")
	if limitStr == "" {
		helpers.HandleError(c, 400, "limit is required")
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request parametr - den page alynyar
	pageStr := c.Query("page")
	if pageStr == "" {
		helpers.HandleError(c, 400, "page is required")
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// request query - den product status alynyar
	// status -> product pozulan ya-da pozulanmadygyny anlatyar
	// true bolsa pozulan
	// false bolsa pozulmadyk
	statusQuery := c.DefaultQuery("status", "false")
	status, err := strconv.ParseBool(statusQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora brend - leryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM products WHERE deleted_at IS NULL`
	if status {
		queryCount = `SELECT COUNT(id) FROM products WHERE deleted_at IS NOT NULL`
	}
	// database - den product - laryn sany alynyar
	var countOfProducts uint
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfProducts); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora product - lary almak ucin query yazylyar
	rowQuery := `SELECT id,name_tm,name_ru,price,old_price,code,brend_id FROM products WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if status {
		rowQuery = `SELECT id,name_tm,name_ru,price,old_price,code,brend_id FROM products WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}

	// database - den brend - lar alynyar
	rowsBrend, err := db.Query(context.Background(), rowQuery, limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsBrend.Close()

	var products []models.Product
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
		); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
		"total":    countOfProducts,
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

	_, err = db.Exec(context.Background(), "UPDATE products SET deleted_at = NOW() WHERE id = $1", ID)
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
	_, err = db.Exec(context.Background(), "UPDATE products SET deleted_at = NULL WHERE id = $1", ID)
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
