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
	HasDelivery bool     `json:"has_delivery,omitempty" binding:"required"`
	ShopOwnerID string   `json:"shop_owner_id,omitempty" binding:"required"`
	SlugTM      string   `json:"slug_tm,omitempty"`
	SlugRU      string   `json:"slug_ru,omitempty"`
	ShopPhones  []string `json:"shop_phones" binding:"required"`
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

func ValidateCreateShop(phoneNumbers []string, orderNumber uint) error {

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

	if orderNumber != 0 {
		var order_number uint
		if err = db.QueryRow(context.Background(), "SELECT order_number FROM shops where order_number = $1", orderNumber).Scan(&order_number); err == nil {
			return errors.New("this order number already exists")
		}
	}

	return nil
}
