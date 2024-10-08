package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/guregu/null.v4"
)

type Shop struct {
	ID               string      `json:"id,omitempty"`
	NameTM           string      `json:"name_tm,omitempty" binding:"required"`
	NameRU           string      `json:"name_ru,omitempty" binding:"required"`
	AddressTM        string      `json:"address_tm,omitempty" binding:"required"`
	AddressRU        string      `json:"address_ru,omitempty" binding:"required"`
	Latitude         float64     `json:"latitude,omitempty"`
	Longitude        float64     `json:"longitude,omitempty"`
	Image            null.String `json:"image,omitempty"`
	HasShipping      bool        `json:"has_shipping,omitempty"`
	ShopOwnerID      null.String `json:"shop_owner_id,omitempty"`
	SlugTM           string      `json:"slug_tm,omitempty"`
	SlugRU           string      `json:"slug_ru,omitempty"`
	ShopPhones       []string    `json:"phones,omitempty"`
	OrderNumber      uint        `json:"order_number,omitempty"`
	IsBrand          bool        `json:"is_brand"`
	IsShoppingCenter bool        `json:"is_shopping_center"`
	ParentShopID     null.String `json:"parent_shop_id,omitempty"`
	CreatedStatus    int8        `json:"created_status,omitempty"`
	AtHome           bool        `json:"at_home"`
}

type UpdateCreatedStatusShop struct {
	ID             string `json:"id" binding:"required"`
	CreatedStatus  int8   `json:"created_status" binding:"required"`
	RejectedReason string `json:"rejected_reason"`
}

type UpdateBrandStatusShop struct {
	ID          string `json:"id" binding:"required"`
	BrandStatus bool   `json:"brand_status"`
}

func ValidateUpdateShopCreatedStatus(shop UpdateCreatedStatusShop) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// body - dan gelen created status dogrymy ya nadogry sol barlanyar
	if shop.CreatedStatus != helpers.CreatedStatuses["wait"] && shop.CreatedStatus != helpers.CreatedStatuses["rejected"] && shop.CreatedStatus != helpers.CreatedStatuses["success"] {
		return errors.New("invalid created status")
	}

	if err := helpers.ValidateRecordByID("shops", shop.ID, "NULL", db); err != nil {
		return err
	}

	if shop.CreatedStatus == helpers.CreatedStatuses["rejected"] && shop.RejectedReason == "" {
		return errors.New(`rejected reason is required`)
	}

	return nil
}

func ValidateCreateShop(shop Shop) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := DefaultValidateShop(shop, db); err != nil {
		return err
	}

	return nil
}

func ValidateUpdateShop(shop Shop) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if shop.ID == "" {
		return errors.New("shop_id is required")
	}
	if err := helpers.ValidateRecordByID("shops", shop.ID, "NULL", db); err != nil {
		return errors.New("record not found")
	}

	if err := DefaultValidateShop(shop, db); err != nil {
		return err
	}

	return nil
}

func DefaultValidateShop(shop Shop, db *pgxpool.Pool) error {
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
		if err := db.QueryRow(
			context.Background(),
			`SELECT id FROM shops WHERE id=$1 AND deleted_at IS NULL AND created_status=$2 AND is_shopping_center=true`, shop.ParentShopID.String, helpers.CreatedStatuses["success"]).
			Scan(&parentShopID); err != nil {
			return errors.New("record not found")
		}

		// egerde shop - yn parenti bar bolsa onda shop awtomat parent shop - ynyn kordinatalaryny almaly
		if err := db.QueryRow(context.Background(), "SELECT latitude,longitude FROM shops WHERE id=$1 AND deleted_at IS NULL", shop.ParentShopID.String).
			Scan(&shop.Latitude, &shop.Longitude); err != nil {
			return err
		}
	}

	if !shop.AtHome && (shop.Latitude == 0 || shop.Longitude == 0) {
		return errors.New("shop coordinates is required")
	}

	return nil
}
