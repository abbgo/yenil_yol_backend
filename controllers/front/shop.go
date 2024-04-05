package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetShops(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// database - den shop - lar alynyar
	rowsShop, err := db.Query(context.Background(), "SELECT id,name_tm,name_ru,latitude,longitude,image,address_tm,address_ru FROM shops WHERE deleted_at IS NULL")
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	var shops []models.Shop
	for rowsShop.Next() {
		var shop models.Shop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude, &shop.Image, &shop.AddressTM, &shop.AddressRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		rowsShopPhones, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id = $1 AND deleted_at IS NULL", shop.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsShopPhones.Close()

		for rowsShopPhones.Next() {
			var phoneNumber string
			if err := rowsShopPhones.Scan(&phoneNumber); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			shop.ShopPhones = append(shop.ShopPhones, phoneNumber)
		}

		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}

func GetShopByIDWithProducts(c *gin.Context) {
	var shop models.Shop
	requestQuery := helpers.StandartQuery{IsDeleted: false}

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

	// request parametrden shop id alynyar
	shopID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,image FROM shops WHERE id = $1 AND deleted_at IS NULL", shopID).Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shop.Image)

	// eger databse sol maglumat yok bolsa error return edilyar
	if shop.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// shop - a degisli telefon belgiler alynyar
	rowsShopPhone, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id=$1 AND deleted_at IS NULL", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShopPhone.Close()
	for rowsShopPhone.Next() {
		var phoneNumber string
		if err := rowsShopPhone.Scan(&phoneNumber); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shop.ShopPhones = append(shop.ShopPhones, phoneNumber)
	}

	// shop - a degisli category - ler alynyar
	rowsCategory, err := db.Query(context.Background(), "SELECT c.id,c.name_tm,c.name_ru FROM categories c INNER JOIN shop_categories sc ON sc.category_id=c.id WHERE sc.shop_id=$1 AND c.parent_category_id IS NULL AND c.deleted_at IS NULL AND sc.deleted_at IS NULL", shop.ID)
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
		shop.ShopCategories = append(shop.ShopCategories, category)
		shop.Categories = append(shop.Categories, category.ID)
	}

	// shop degisli product - lar alynyar
	rowsProducts, err := db.Query(context.Background(), "SELECT DISTINCT ON (p.id,p.created_at) p.id,p.name_tm,p.name_ru,p.price,p.old_price FROM products p INNER JOIN category_products cp ON cp.product_id=p.id WHERE cp.category_id=ANY($1) AND cp.deleted_at IS NULL AND p.deleted_at IS NULL ORDER BY p.created_at DESC LIMIT $2 OFFSET $3", pq.Array(shop.Categories), requestQuery.Limit, offset)
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

		shop.CountOfProducts++
		shop.Products = append(shop.Products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})
}
