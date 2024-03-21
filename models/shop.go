package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type Shop struct {
	ID          string   `json:"id,omitempty"`
	NameTM      string   `json:"name_tm,omitempty" binding:"required"`
	NameRU      string   `json:"name_ru,omitempty" binding:"required"`
	AddressTM   string   `json:"address_tm,omitempty" binding:"required"`
	AddressRU   string   `json:"address_ru,omitempty" binding:"required"`
	Latitude    float64  `json:"latitude,omitempty" binding:"required"`
	Longitude   float64  `json:"longitude,omitempty" binding:"required"`
	Image       string   `json:"image,omitempty"`
	HasDelivery bool     `json:"has_delivery,omitempty"`
	ShopOwnerID string   `json:"shop_owner_id,omitempty" binding:"required"`
	SlugTM      string   `json:"slug_tm,omitempty"`
	SlugRU      string   `json:"slug_ru,omitempty"`
	ShopPhones  []string `json:"phones" binding:"required"`
	Categories  []string `json:"categories" binding:"required"`
	OrderNumber uint     `json:"order_number,omitempty"`
	CreatedAt   string   `json:"-"`
	UpdatedAt   string   `json:"-"`
	DeletedAt   string   `json:"-"`
}

type ShopQuery struct {
	helpers.StandartQuery
	ShopOwnerID string `form:"shop_owner_id"`
}

func ValidateShop(shop Shop, isCreateFunction bool) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// telefon belgiler barlanylyar
	for _, v := range shop.ShopPhones {
		if !helpers.ValidatePhoneNumber(v) {
			return errors.New("invalid phone number")
		}
	}

	if err := helpers.ValidateRecordByID("shop_owners", shop.ShopOwnerID, "NULL", db); err != nil {
		return err
	}

	for _, v := range shop.Categories {
		if err := helpers.ValidateRecordByID("categories", v, "NULL", db); err != nil {
			return err
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
			var shop_id string
			db.QueryRow(context.Background(), "SELECT id FROM shops where order_number = $1 AND deleted_at IS NULL", shop.OrderNumber).Scan(&shop_id)
			if shop_id != shop.ID && shop_id != "" {
				return errors.New("this order number already exists")
			}
		}
	}

	return nil
}
