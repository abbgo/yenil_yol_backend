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

// bu model ShopOwner - in maglumatyny uytgetmegi ucin doredildi
type ShopOwnerUpdate struct {
	ID          string `json:"id,omitempty" binding:"required"`
	FullName    string `json:"full_name,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
}

func ValidateShopOwner(phoneNumber, shopOwnerID string, isRegisterFunction bool) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if isRegisterFunction {
		var phone_number string
		db.QueryRow(context.Background(), "SELECT phone_number FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&phone_number)
		if phone_number != "" {
			return errors.New("this shop owner already exists")
		}
	} else {
		if shopOwnerID == "" {
			return errors.New("shop_owner_id is required")
		}

		if err := helpers.ValidateRecordByID("shop_owners", shopOwnerID, "NULL", db); err != nil {
			return err
		}

		var shop_owner_id string
		db.QueryRow(context.Background(), "SELECT id FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", phoneNumber).Scan(&shop_owner_id)
		if shop_owner_id != shopOwnerID && shop_owner_id != "" {
			return errors.New("this shop owner already exists")
		}
	}

	if !helpers.ValidatePhoneNumber(phoneNumber) {
		return errors.New("invalid phone number")
	}

	return nil
}
