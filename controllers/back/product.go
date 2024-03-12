package controllers

import (
	"context"
	"database/sql"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
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

	if err := models.ValidateProduct(product.Price, product.OldPrice); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda products tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO products (name_tm,name_ru,image,price,old_price,status,color_name_tm,color_name_ru,gender_name_tm,gender_name_ru,code,slug_tm,slug_ru,shop_id,category_id,brend_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)", product.NameTM, product.NameRU, product.Image, product.Price, product.OldPrice, product.Status, product.ColorNameTM, product.ColorNameRU, product.GenderNameTM, product.GenderNameRU, product.Code, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.ShopID, product.CategoryID, product.BrendID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa brend ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", product.Image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
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
	var product models.ProductUpdate
	if err := c.BindJSON(&product); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateProduct(product.Price, product.OldPrice); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var productdID string
	var oldProductImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM products WHERE id = $1 AND deleted_at IS NULL", product.ID).Scan(&productdID, &oldProductImage); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if productdID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// bu yerde product - yn suraty ucin fileName atly uytgeyan ululyk doredilyar
	// eger request body - dan surat gelmese onda product - yn suraty uytgedilmeyar diymek bolyar
	// sonun ucin product - in onki suratyny goyyas , eger request body - dan surat gelen bolsa
	// onda taze suraty kone surat bilen calysyas
	var fileName string
	if product.Image == "" {
		fileName = oldProductImage.String
	} else {
		// sonra helper_images tablisa product ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", product.Image)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if oldProductImage.String != "" {
			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + oldProductImage.String); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		fileName = product.Image
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE products SET name_tm=$1 , name_ru=$2 , image=$3 , price=$4 , old_price=$5 , status=$6 , color_name_tm=$7 , color_name_ru=$8 , gender_name_tm=$9 , gender_name_ru=$10 , code=$11 , slug_tm=$12 , slug_ru=$13 , shop_id=$14 , category_id=$15 , brend_id=$16 WHERE id=$17", product.NameTM, product.NameRU, fileName, product.Price, product.OldPrice, product.Status, product.ColorNameTM, product.ColorNameRU, product.GenderNameTM, product.GenderNameRU, product.Code, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.ShopID, product.CategoryID, product.BrendID, product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
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
	var productImage sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,image,price,old_price,status,color_name_tm,color_name_ru,gender_name_tm,gender_name_ru,code,shop_id,category_id,brend_id FROM products WHERE id = $1 AND deleted_at IS NULL", productID).Scan(
		&product.ID,
		&product.NameTM,
		&product.NameRU,
		&productImage,
		&product.Price,
		&product.OldPrice,
		&product.Status,
		&product.ColorNameTM,
		&product.ColorNameRU,
		&product.GenderNameTM,
		&product.GenderNameRU,
		&product.Code,
		&product.ShopID,
		&product.CategoryID,
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

	if productImage.String != "" {
		product.Image = productImage.String
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
	rowQuery := `SELECT id,name_tm,name_ru,image,price,old_price,status,color_name_tm,color_name_ru,gender_name_tm,gender_name_ru,code,shop_id,category_id,brend_id FROM products WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if status {
		rowQuery = `SELECT id,name_tm,name_ru,image,price,old_price,status,color_name_tm,color_name_ru,gender_name_tm,gender_name_ru,code,shop_id,category_id,brend_id FROM products WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
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
		var productImage sql.NullString
		if err := rowsBrend.Scan(
			&product.ID,
			&product.NameTM,
			&product.NameRU,
			&productImage,
			&product.Price,
			&product.OldPrice,
			&product.Status,
			&product.ColorNameTM,
			&product.ColorNameRU,
			&product.GenderNameTM,
			&product.GenderNameRU,
			&product.Code,
			&product.ShopID,
			&product.CategoryID,
			&product.BrendID,
		); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if productImage.String != "" {
			product.Image = productImage.String
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

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM products WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa shop we sol brend - in deleted_at - ine current_time goyulyar
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

	// alynan id den bolan product database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM products WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database sol id degisli product yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
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

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	var image sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT id,image FROM products WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id, &image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli product yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger shop bar bolsa sonda product - yn suraty papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
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
