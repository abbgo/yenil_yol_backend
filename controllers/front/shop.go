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
	ID         string   `json:"id,omitempty"`
	NameTM     string   `json:"name_tm,omitempty"`
	NameRU     string   `json:"name_ru,omitempty"`
	Latitude   float64  `json:"latitude,omitempty"`
	Longitude  float64  `json:"longitude,omitempty"`
	Image      string   `json:"image,omitempty"`
	ShopPhones []string `json:"shop_phones"`
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
	rowsShop, err := db.Query(context.Background(), "SELECT id,name_tm,name_ru,latitude,longitude,image FROM shops WHERE deleted_at IS NULL")
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	var shops []Shop
	for rowsShop.Next() {
		var shop Shop
		var shopImage sql.NullString
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude, &shopImage); err != nil {
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
