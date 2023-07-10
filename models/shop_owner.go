package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type ShopOwner struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password,omitempty" binding:"required"`
	Slug        string `json:"slug,omitempty"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}

// bu model ShopOwner Login bolmagy ucin doredildi
type ShopOwnerLogin struct {
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password,omitempty" binding:"required"`
}

// bu model ShopOwner - in maglumatyny uytgetmegi ucin doredildi
type ShopOwnerUpdate struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
}

func ValidateRegisterShopOwner(phoneNumber, url string) error {

	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if url == "register" {
		var phone_number string
		db.QueryRow(context.Background(), "SELECT phone_number FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&phone_number)
		if phone_number != "" {
			return errors.New("this shop owner already exists")
		}
	}

	if !helpers.ValidatePhoneNumber(phoneNumber) {
		return errors.New("invalid phone number")
	}

	return nil

}
