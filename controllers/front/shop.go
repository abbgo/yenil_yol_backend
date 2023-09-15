package controllers

import (
	"context"
	"database/sql"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Shop struct {
	ID          string   `json:"id,omitempty"`
	NameTM      string   `json:"name_tm,omitempty"`
	NameRU      string   `json:"name_ru,omitempty"`
	Latitude    float64  `json:"latitude,omitempty"`
	Longitude   float64  `json:"longitude,omitempty"`
	Image       string   `json:"image,omitempty"`
	ShopPhones  []string `json:"shop_phones"`
	AddressTM   string   `json:"address_tm,omitempty"`
	AddressRU   string   `json:"address_ru,omitempty"`
	HasDelivery bool     `json:"has_delivery,omitempty"`
	ShopOwnerID string   `json:"shop_owner_id,omitempty"`
}

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

	var shops []Shop
	for rowsShop.Next() {
		var shop Shop
		var shopImage sql.NullString
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude, &shopImage, &shop.AddressTM, &shop.AddressRU); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if shopImage.String != "" {
			shop.Image = shopImage.String
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
	rowShop, err := db.Query(context.Background(), "SELECT s.id,s.name_tm,s.name_ru,s.address_tm,s.address_ru,s.latitude,s.longitude,s.image,sp.phone_number FROM shops s INNER JOIN shop_phones sp ON sp.shop_id = s.id WHERE s.id = $1 AND s.deleted_at IS NULL AND sp.deleted_at IS NULL", shopID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowShop.Close()

	var shop Shop
	var shopImage sql.NullString
	for rowShop.Next() {
		var shopPhone string
		if err := rowShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shopImage, &shopPhone); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		if shopImage.String != "" {
			shop.Image = shopImage.String
		}
		shop.ShopPhones = append(shop.ShopPhones, shopPhone)
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if shop.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})

}
