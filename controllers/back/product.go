package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"
	"os"

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
	if err := db.QueryRow(context.Background(),
		`INSERT INTO products (name_tm,name_ru,price,old_price,code,slug_tm,slug_ru,brend_id,is_visible,shop_id) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id`,
		product.NameTM, product.NameRU, product.Price, product.OldPrice, productCode, slug.MakeLang(product.NameTM, "en"),
		slug.MakeLang(product.NameRU, "en"), product.BrendID, product.IsVisible, product.ShopID).Scan(&product.ID); err != nil {
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
		if err := db.QueryRow(context.Background(),
			`INSERT INTO product_colors (name,product_id,order_number) VALUES ($1,$2,$3) RETURNING id`,
			v.Name, product.ID, v.OrderNumber).Scan(&productColorID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// renk gosulandan sonra sol renke degisli razmerler gosulyar
		_, err := db.Exec(context.Background(),
			`INSERT INTO product_dimensions (product_color_id,dimension_id) VALUES ($1,unnest($2::uuid[]))`,
			productColorID, pq.Array(v.Dimensions))
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		for _, image := range v.Images {
			// bu yerde renke degisli suratlar gosulyar
			_, err = db.Exec(context.Background(),
				"INSERT INTO product_images (product_color_id,image,resized_image,order_number) VALUES ($1,$2,$3,$4)",
				productColorID, image.Image, "assets/"+image.Image, image.OrderNumber)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
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

		for _, image := range v.Images {
			// resizedImages = append(resizedImages, "assets/"+rm)
			// bu yerde renke degisli suratlar gosulyar
			_, err = db.Exec(context.Background(),
				"INSERT INTO product_images (product_color_id,image,resized_image,order_number) VALUES ($1,$2,$3,$4)",
				productColorID, image.Image, "assets/"+image.Image, image.OrderNumber)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
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
	var product serializations.GetProductForBack
	if err := db.QueryRow(context.Background(),
		`SELECT id,name_tm,name_ru,price,old_price,brend_id,is_visible FROM products WHERE id = $1 AND deleted_at IS NULL`,
		productID).
		Scan(
			&product.ID,
			&product.NameTM,
			&product.NameRU,
			&product.Price,
			&product.OldPrice,
			&product.BrendID,
			&product.IsVisible,
		); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if product.BrendID.String != "" {
		if err := db.QueryRow(context.Background(), `SELECT id,name FROM brends WHERE id=$1`, product.BrendID.String).
			Scan(&product.Brend.ID, &product.Brend.Name); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// haryda degisli category - lar alynyar
	rowsCategory, err := db.Query(
		context.Background(),
		`SELECT c.id,c.name_tm,c.name_ru FROM categories c 
		INNER JOIN category_products cp ON cp.category_id=c.id 
		WHERE cp.product_id=$1 AND cp.deleted_at IS NULL AND c.deleted_at IS NULL`,
		productID)
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
		product.Categories = append(product.Categories, category)
	}

	// haryda degisli renkler , razmerler we suratlar alynyar
	rowsColor, err := db.Query(context.Background(),
		`SELECT id,name,order_number FROM product_colors WHERE product_id=$1 AND deleted_at IS NULL`,
		productID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsColor.Close()

	for rowsColor.Next() {
		var productColor serializations.ProductColorForBack
		if err := rowsColor.Scan(&productColor.ID, &productColor.Name, &productColor.OrderNumber); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// renk alynandan son sol renke degisli razmerler alynyar
		rowsDimension, err := db.Query(context.Background(),
			`SELECT d.id,d.dimension FROM dimensions d 
			INNER JOIN product_dimensions pd ON pd.dimension_id=d.id 
			WHERE pd.product_color_id=$1 
			AND pd.deleted_at IS NULL AND d.deleted_at IS NULL`,
			productColor.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsDimension.Close()

		for rowsDimension.Next() {
			var dimension models.Dimension
			if err := rowsDimension.Scan(&dimension.ID, &dimension.Dimension); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			productColor.Dimensions = append(productColor.Dimensions, dimension)
		}

		// renk alynandan son sol renke degisli suratlar alynyar
		rowsImage, err := db.Query(context.Background(),
			`SELECT image,order_number FROM product_images WHERE product_color_id=$1 AND deleted_at IS NULL`,
			productColor.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsImage.Close()

		for rowsImage.Next() {
			var image models.ProductImage
			if err := rowsImage.Scan(&image.Image, &image.OrderNumber); err != nil {
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
	var requestQuery serializations.ProductQuery
	var products []serializations.GetProductsForBack
	isDeleted := "NULL"
	var shopQuery string

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

	// request query - den status - a gora product - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(`SELECT p.id,p.name_tm,p.name_ru,p.price,p.old_price,p.is_visible FROM products p WHERE p.deleted_at IS %v`, isDeleted)
	orderQuery := fmt.Sprintf(` ORDER BY p.created_at DESC LIMIT %v OFFSET %v`, requestQuery.Limit, offset)
	if requestQuery.ShopID != "" {
		shopQuery = fmt.Sprintf(` AND p.shop_id = '%v'`, requestQuery.ShopID)
	}

	// database - den brend - lar alynyar
	rowsProducts, err := db.Query(context.Background(), rowQuery+shopQuery+orderQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsProducts.Close()

	for rowsProducts.Next() {
		var product serializations.GetProductsForBack
		if err := rowsProducts.Scan(
			&product.ID,
			&product.NameTM,
			&product.NameRU,
			&product.Price,
			&product.OldPrice,
			&product.IsVisible,
		); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := db.QueryRow(context.Background(),
			`SELECT image FROM product_images pi INNER JOIN product_colors pc ON pc.id=pi.product_color_id WHERE pc.product_id=$1 AND pc.order_number=1 AND pi.order_number=1`,
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
