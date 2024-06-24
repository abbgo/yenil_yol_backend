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
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

func GetShopsForMap(c *gin.Context) {
	var requestQuery models.ShopForMapQuery

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

	rowsShopQuery := fmt.Sprintf(`
							SELECT id,name_tm,name_ru,latitude,longitude FROM shops 
							WHERE 6371 * acos(
										cos(radians(%f)) * cos(radians(latitude)) *
										cos(radians(longitude) - radians(%f)) +
										sin(radians(%f)) * sin(radians(latitude))
									) <= %d AND deleted_at IS NULL;
						`, requestQuery.Latitude, requestQuery.Longitude, requestQuery.Latitude, requestQuery.Kilometer)

	// if requestQuery.Search != "" {
	// 	rowsShopQuery = fmt.Sprintf(`
	// 						SELECT id,name_tm,name_ru,latitude,longitude FROM shops
	// 						WHERE 6371 * acos(
	// 									cos(radians(%f)) * cos(radians(latitude)) *
	// 									cos(radians(longitude) - radians(%f)) +
	// 									sin(radians(%f)) * sin(radians(latitude))
	// 								) <= %d AND to_tsvector(slug_tm) @@ to_tsquery('%s') OR slug_tm LIKE '%s'  AND deleted_at IS NULL;
	// 					`, requestQuery.Latitude, requestQuery.Longitude, requestQuery.Latitude, requestQuery.Kilometer, search, searchStr)
	// }

	rowsShop, err := db.Query(context.Background(), rowsShopQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()
	var shops []models.Shop
	for rowsShop.Next() {
		var shop models.Shop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}

func GetShops(c *gin.Context) {
	var requestQuery models.ShopQuery
	var search, searchStr, queryRandom, querySearch string

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

	// database - den shop - lar alynyar
	queryDefault := `SELECT id,name_tm,name_ru,latitude,longitude,resized_image,address_tm,address_ru,is_brend FROM shops WHERE deleted_at IS NULL`
	if requestQuery.IsRandom {
		queryRandom = ` ORDER BY RANDOM()`
	}
	if requestQuery.Search != "" {
		querySearch = fmt.Sprintf(` AND (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s')`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}
	queryLimitOffset := fmt.Sprintf(` LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	rowsShop, err := db.Query(context.Background(), queryDefault+querySearch+queryRandom+queryLimitOffset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	var shops []models.Shop
	for rowsShop.Next() {
		var shop models.Shop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude, &shop.Image, &shop.AddressTM, &shop.AddressRU, &shop.IsBrend); err != nil {
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

func GetShopByID(c *gin.Context) {
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
	db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,resized_image,has_delivery FROM shops WHERE id = $1 AND deleted_at IS NULL", shopID).Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shop.Image, &shop.HasDelivery)

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

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})
}

func GetShopByIDs(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id - ler alynyar
	shopIDs := c.QueryArray("ids")

	// database - den request parametr - den gelen id - ler boyunca maglumat cekilyar
	var shops []models.Shop
	rows, err := db.Query(context.Background(),
		`
			SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,resized_image FROM shops 
			WHERE id = ANY($1) AND deleted_at IS NULL
		`,
		pq.Array(shopIDs))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var shop models.Shop
		if err := rows.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shop.Image); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}
