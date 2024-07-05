package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"

	"gopkg.in/guregu/null.v4"
)

type Shop struct {
	ID               string      `json:"id,omitempty"`
	NameTM           string      `json:"name_tm,omitempty" binding:"required"`
	NameRU           string      `json:"name_ru,omitempty" binding:"required"`
	AddressTM        string      `json:"address_tm,omitempty" binding:"required"`
	AddressRU        string      `json:"address_ru,omitempty" binding:"required"`
	Latitude         float64     `json:"latitude,omitempty" binding:"required"`
	Longitude        float64     `json:"longitude,omitempty" binding:"required"`
	Image            null.String `json:"image,omitempty"`
	HasShipping      bool        `json:"has_shipping,omitempty"`
	ShopOwnerID      null.String `json:"shop_owner_id,omitempty"`
	SlugTM           string      `json:"slug_tm,omitempty"`
	SlugRU           string      `json:"slug_ru,omitempty"`
	ShopPhones       []string    `json:"phones,omitempty"`
	OrderNumber      uint        `json:"order_number,omitempty"`
	IsBrend          bool        `json:"is_brend"`
	IsShoppingCenter bool        `json:"is_shopping_center"`
	ParentShopID     null.String `json:"parent_shop_id,omitempty"`
}

type ShopQuery struct {
	helpers.StandartQuery
	ShopOwnerID      string `form:"shop_owner_id"`
	IsRandom         bool   `form:"is_random"`
	Search           string `form:"search"`
	Lang             string `form:"lang"`
	IsShoppingCenter bool   `form:"is_shopping_center"`
}

type ShopForMapQuery struct {
	Latitude  float64 `form:"latitude" validate:"required"`
	Longitude float64 `form:"longitude" validate:"required"`
	Kilometer int8    `form:"kilometer"`
}

func ValidateShop(shop Shop, isCreateFunction bool) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !shop.IsShoppingCenter && shop.ShopOwnerID.String == "" {
		return errors.New("shop_owner_id is required")
	}

	if !shop.IsShoppingCenter && len(shop.ShopPhones) == 0 {
		return errors.New("shop_phones is required")
	}

	// telefon belgiler barlanylyar
	if len(shop.ShopPhones) != 0 {
		for _, v := range shop.ShopPhones {
			if !helpers.ValidatePhoneNumber(v) {
				return errors.New("invalid phone number")
			}
		}
	}

	if shop.ShopOwnerID.String != "" {
		if err := helpers.ValidateRecordByID("shop_owners", shop.ShopOwnerID.String, "NULL", db); err != nil {
			return err
		}
	}

	if shop.ParentShopID.String != "" {
		var parentShopID string
		if err := db.QueryRow(context.Background(), "SELECT id FROM shops WHERE deleted_at IS NULL AND is_shopping_center=true").Scan(&parentShopID); err != nil {
			return errors.New("record not found")
		}

		// egerde shop - yn parenti bar bolsa onda shop awtomat parent shop - ynyn kordinatalaryny almaly
		if err := db.QueryRow(context.Background(), "SELECT latitude,longitude FROM shops WHERE id=$1 AND deleted_at IS NULL", shop.ParentShopID.String).Scan(&shop.Latitude, &shop.Longitude); err != nil {
			return err
		}
	}

	if !isCreateFunction {
		if shop.ID == "" {
			return errors.New("shop_id is required")
		}

		if err := helpers.ValidateRecordByID("shops", shop.ID, "NULL", db); err != nil {
			return errors.New("record not found")
		}
	}

	if shop.OrderNumber != 0 {
		if isCreateFunction {
			var order_number uint
			if err = db.QueryRow(context.Background(), "SELECT order_number FROM shops where order_number = $1 AND deleted_at IS NULL AND shop_owner_id = $2", shop.OrderNumber, shop.ShopOwnerID).Scan(&order_number); err == nil {
				return errors.New("this order number already exists")
			}
		} else {
			if shop.ID == "" {
				return errors.New("shop_id is required")
			}

			if err := helpers.ValidateRecordByID("shops", shop.ID, "NULL", db); err != nil {
				return errors.New("record not found")
			}

			var shop_id string
			db.QueryRow(context.Background(), "SELECT id FROM shops where order_number = $1 AND deleted_at IS NULL", shop.OrderNumber).Scan(&shop_id)
			if shop_id != shop.ID && shop_id != "" {
				return errors.New("this order number already exists")
			}
		}
	}

	return nil
}
