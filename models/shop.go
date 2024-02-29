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
	OrderNumber uint     `json:"order_number,omitempty"`
	CreatedAt   string   `json:"-"`
	UpdatedAt   string   `json:"-"`
	DeletedAt   string   `json:"-"`
}

type ShopPhone struct {
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	ShopID      string `json:"shop_id,omitempty" binding:"required"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

func ValidateShop(phoneNumbers []string, orderNumber uint, isCreateFunction bool, shopId, shopOwnerID string) error {

	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// telefon belgiler barlanylyar
	for _, v := range phoneNumbers {
		if !helpers.ValidatePhoneNumber(v) {
			return errors.New("invalid phone number")
		}
	}

	if shopOwnerID == "" {
		return errors.New("shop_owner_id is required")
	}
	var shop_owner_id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM shop_owners WHERE id = $1 AND deleted_at IS NULL", shopOwnerID).Scan(&shop_owner_id); err != nil {
		return err
	}
	if shop_owner_id == "" {
		return errors.New("shop_owner not found")
	}

	if orderNumber != 0 {
		if isCreateFunction {
			var order_number uint
			if err = db.QueryRow(context.Background(), "SELECT order_number FROM shops where order_number = $1 AND deleted_at IS NULL AND shop_owner_id = $2", orderNumber, shopOwnerID).Scan(&order_number); err == nil {
				return errors.New("this order number already exists")
			}
		} else {
			if shopId == "" {
				return errors.New("shop_id is required")
			}
			var shop_id string
			db.QueryRow(context.Background(), "SELECT id FROM shops where order_number = $1 AND deleted_at IS NULL", orderNumber).Scan(&shop_id)
			if shop_id != shopId && shop_id != "" {
				return errors.New("this order number already exists")
			}
		}
	}

	return nil
}
