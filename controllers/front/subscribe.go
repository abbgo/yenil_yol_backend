package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShopID struct {
	IDS []string `json:"shop_ids"`
}

func AddOrRemoveSubscribe(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// bu yerde api query - den status ady bilen status alynyar
	// eger status = true bolsa onda halanlarym sahypa haryt gosulyar
	// eger status = false bolsa onda halanlarym sahypadan haryt ayrylyar
	statusStr := c.DefaultQuery("status", "true")
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// bu yerde middleware - den gelen musderinin id - si alynyar
	custID, hasCustomer := c.Get("customer_id")
	if !hasCustomer {
		helpers.HandleError(c, 400, "customer_id is required")
		return
	}
	customerID, ok := custID.(string)
	if !ok {
		helpers.HandleError(c, 400, "customer_id must be string")
		return
	}

	// front - dan gelen maglumatlar bind edilyar
	var shopIDs ShopID
	if err := c.BindJSON(&shopIDs); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if status { // eger status = true bolsa subscribe sahypa magazyn gosulyar
		if len(shopIDs.IDS) != 0 { // eger front - dan gelen id bar bolsa onda subscribe sahypa haryt gosulyar
			for _, v := range shopIDs.IDS {
				// front - dan gelen id - lere den gelyan bazada magazyn barmy yokmy sol barlanyar
				var shop_id string
				if err := db.QueryRow(context.Background(), "SELECT id FROM shops WHERE id = $1 AND deleted_at IS NULL", v).Scan(&shop_id); err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}

				if shop_id != "" { // eger magazyn shops tablida bar bolsa onda sol magazyn on gelen musderinin
					// subscribe - larynyn arasynda barmy ya-da yokmy sol barlanyar
					var shop string
					db.QueryRow(context.Background(), "SELECT shop_id FROM subscribes WHERE customer_id = $1 AND shop_id = $2 AND deleted_at IS NULL", customerID, v).Scan(&shop)

					if shop == "" { // eger magazyn musderinin subscribe - larynyn arasynda yok bolsa
						// onda magazyn sol musderinin subscribes tablisasyna gosulyar
						_, err := db.Exec(context.Background(), "INSERT INTO subscribes (customer_id,shop_id) VALUES ($1,$2)", customerID, v)
						if err != nil {
							helpers.HandleError(c, 400, err.Error())
							return
						}
					}
				}
			}

			// front - dan gelen magazynlaryn subscribe sahypa gosulandan son
			// yzyna sol magazynlar ddoly maglumatlary bilen berilyar
			shops, err := GetSubscribes(customerID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"shops":  shops,
			})
		} else { // eger front hic hile id gelmese onda musderinin onki bazadaky subscribe magazynlaryny fronta bermeli
			shops, err := GetSubscribes(customerID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			if len(shops) == 0 {
				helpers.HandleError(c, 400, "subscribe empty")
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status": true,
				"shops":  shops,
			})
		}
	} else { // eger status = false gelse onda front - dan gele id - li magazyny sol musderinin subscribe magazynlaryndan pozmaly
		if len(shopIDs.IDS) != 0 { // front - dan maglumat gelyarmi ya-da gelenokmy sony barlayas
			// front - dan gelen id - ler subscribes tablisada barmy ya-da yokmy sony barlayas
			var shop_id string
			db.QueryRow(context.Background(), "SELECT shop_id FROM subscribes WHERE customer_id = $1 AND shop_id = $2 AND deleted_at IS NULL", customerID, shopIDs.IDS[0]).Scan(&shop_id)

			// eger magazyn subscribes tablisada yok bolsa
			// yzyna yalnyslyk goyberyas
			if shop_id == "" {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": "this shop not found in this customer",
				})
				return
			}

			// magazyn subscribe tablisada bar bolsa onda ony pozyas
			_, err = db.Exec(context.Background(), "DELETE FROM subscribes WHERE customer_id = $1 AND shop_id = $2 AND deleted_at IS NULL", customerID, shopIDs.IDS[0])
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "subscribe successfull deleted",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "shop id is required",
			})
		}
	}
}

func GetSubscribes(customerID string) ([]models.Shop, error) {

	db, err := config.ConnDB()
	if err != nil {
		return []models.Shop{}, err
	}
	defer db.Close()

	rowsShop, err := db.Query(context.Background(),
		"SELECT sh.id,sh.name_tm,sh.name_ru,sh.address_tm,sh.address_ru,sh.latitude,sh.longitude,sh.image,sh.has_delivery,sh.shop_owner_id,sh.order_number FROM shops sh INNER JOIN subscribes su ON su.shop_id = sh.id WHERE su.customer_id = $1 AND su.deleted_at IS NULL AND sh.deleted_at IS NULL",
		customerID)
	if err != nil {
		return []models.Shop{}, err
	}
	defer rowsShop.Close()

	var shops []models.Shop
	for rowsShop.Next() {
		var shop models.Shop
		if err := rowsShop.Scan(
			&shop.ID,
			&shop.NameTM,
			&shop.NameRU,
			&shop.AddressTM,
			&shop.AddressRU,
			&shop.Latitude,
			&shop.Longitude,
			&shop.Image,
			&shop.HasDelivery,
			&shop.ShopOwnerID,
			&shop.OrderNumber,
		); err != nil {
			return []models.Shop{}, err
		}

		rowsShopPhone, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id = $1 AND deleted_at IS NULL", shop.ID)
		if err != nil {
			return []models.Shop{}, err
		}
		defer rowsShopPhone.Close()

		for rowsShopPhone.Next() {
			var shopPhone string
			if err := rowsShopPhone.Scan(&shopPhone); err != nil {
				return []models.Shop{}, err
			}

			shop.ShopPhones = append(shop.ShopPhones, shopPhone)
		}

		shops = append(shops, shop)
	}

	return shops, nil

}

func GetCustomerSubscribes(c *gin.Context) {

	custID, hasCustomer := c.Get("customer_id")
	if !hasCustomer {
		helpers.HandleError(c, 400, "customer_id is required")
		return
	}
	customerID, ok := custID.(string)
	if !ok {
		helpers.HandleError(c, 400, "customer_id must be string")
		return
	}

	shops, err := GetSubscribes(customerID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})

}
