package serializations

import (
	"github/abbgo/yenil_yol/backend/helpers"

	"gopkg.in/guregu/null.v4"
)

type ShopQuery struct {
	helpers.StandartQuery
	ShopOwnerID      string   `form:"shop_owner_id"`
	IsRandom         bool     `form:"is_random"`
	Search           string   `form:"search"`
	Lang             string   `form:"lang"`
	IsShoppingCenter bool     `form:"is_shopping_center"`
	ParentShopID     string   `form:"parent_shop_id"`
	CratedStatuses   []string `form:"crated_statuses"`
}

type ShopForMapQuery struct {
	Latitude  float64  `form:"latitude" validate:"required"`
	Longitude float64  `form:"longitude" validate:"required"`
	Kilometer int8     `form:"kilometer"`
	Genders   []string `form:"genders"`
}

type GetShops struct {
	ID             string      `json:"id,omitempty"`
	NameTM         string      `json:"name_tm,omitempty"`
	NameRU         string      `json:"name_ru,omitempty"`
	Latitude       float64     `json:"latitude,omitempty"`
	Longitude      float64     `json:"longitude,omitempty"`
	Image          null.String `json:"image,omitempty"`
	HasShipping    bool        `json:"has_shipping"`
	CreatedStatus  int8        `json:"created_status"`
	RejectedReason null.String `json:"rejected_reason,omitempty"`
}

type GetShop struct {
	ID               string      `json:"id,omitempty"`
	NameTM           string      `json:"name_tm,omitempty"`
	NameRU           string      `json:"name_ru,omitempty"`
	AddressTM        string      `json:"address_tm,omitempty"`
	AddressRU        string      `json:"address_ru,omitempty"`
	Latitude         float64     `json:"latitude,omitempty"`
	Longitude        float64     `json:"longitude,omitempty"`
	Image            null.String `json:"image,omitempty"`
	HasShipping      bool        `json:"has_shipping,omitempty"`
	ShopOwnerID      null.String `json:"shop_owner_id,omitempty"`
	ShopOwner        *ShopOwner  `json:"shop_owner,omitempty"`
	ShopPhones       []string    `json:"phones,omitempty"`
	ParentShopID     null.String `json:"-"`
	ParentShop       *ParentShop `json:"parent_shop"`
	IsShoppingCenter bool        `json:"is_shopping_center,omitempty"`
	AtHome           bool        `json:"at_home"`
	IsBrand          bool        `json:"is_brand"`
}

type ParentShop struct {
	ID               string `json:"id,omitempty"`
	NameTM           string `json:"name_tm,omitempty"`
	NameRU           string `json:"name_ru,omitempty"`
	IsShoppingCenter bool   `json:"is_shopping_center,omitempty"`
}
