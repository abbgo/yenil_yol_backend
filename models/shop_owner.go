package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
)

type ShopOwner struct {
	ID          string `json:"id,omitempty"`
	FullName    string `json:"full_name,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password,omitempty" binding:"required"`
	Slug        string `json:"slug,omitempty"`
}

// bu model ShopOwner Login bolmagy ucin doredildi
type ShopOwnerLogin struct {
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password,omitempty" binding:"required"`
}

func ValidateShopOwner(shopOwner ShopOwner, isRegisterFunction bool) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if isRegisterFunction {
		var phone_number string
		db.QueryRow(context.Background(), "SELECT phone_number FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", shopOwner.PhoneNumber).Scan(&phone_number)
		if phone_number != "" {
			return errors.New("this shop owner already exists")
		}
	} else {
		if shopOwner.ID == "" {
			return errors.New("shop_owner_id is required")
		}

		if err := helpers.ValidateRecordByID("shop_owners", shopOwner.ID, "NULL", db); err != nil {
			return err
		}

		var shop_owner_id string
		db.QueryRow(context.Background(), "SELECT id FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", shopOwner.PhoneNumber).Scan(&shop_owner_id)
		if shop_owner_id != shopOwner.ID && shop_owner_id != "" {
			return errors.New("this shop owner already exists")
		}
	}

	if !helpers.ValidatePhoneNumber(shopOwner.PhoneNumber) {
		return errors.New("invalid phone number")
	}

	return nil
}
