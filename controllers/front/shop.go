package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
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
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	shopID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var shop models.Shop
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
	}

	// rowsProducts, err := db.Query(context.Background(), "SELECT id,name_tm,name_ru,price,old_price FROM products WHERE shop_id = $1 AND deleted_at IS NULL", shop.ID)
	// if err != nil {
	// 	helpers.HandleError(c, 400, err.Error())
	// 	return
	// }

	// for rowsProducts.Next() {
	// 	var product models.Product
	// 	if err := rowsProducts.Scan(&product.ID, &product.NameTM, &product.NameRU, &product.Price, &product.OldPrice); err != nil {
	// 		helpers.HandleError(c, 400, err.Error())
	// 		return
	// 	}
	// 	shop.Products = append(shop.Products, product)
	// }

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})
}
