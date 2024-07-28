package serializations

import (
	"github/abbgo/yenil_yol/backend/helpers"

	"gopkg.in/guregu/null.v4"
)

type ShopQuery struct {
	helpers.StandartQuery
	ShopOwnerID      string `form:"shop_owner_id"`
	IsRandom         bool   `form:"is_random"`
	Search           string `form:"search"`
	Lang             string `form:"lang"`
	IsShoppingCenter bool   `form:"is_shopping_center"`
	ParentShopID     string `form:"parent_shop_id"`
}

type ShopForMapQuery struct {
	Latitude  float64 `form:"latitude" validate:"required"`
	Longitude float64 `form:"longitude" validate:"required"`
	Kilometer int8    `form:"kilometer"`
}

type GetShops struct {
	ID          string      `json:"id,omitempty"`
	NameTM      string      `json:"name_tm,omitempty"`
	NameRU      string      `json:"name_ru,omitempty"`
	Latitude    float64     `json:"latitude,omitempty"`
	Longitude   float64     `json:"longitude,omitempty"`
	Image       null.String `json:"image,omitempty"`
	HasShipping bool        `json:"has_shipping"`
}

type GetShop struct {
	ID           string      `json:"id,omitempty"`
	NameTM       string      `json:"name_tm,omitempty"`
	NameRU       string      `json:"name_ru,omitempty"`
	AddressTM    string      `json:"address_tm,omitempty"`
	AddressRU    string      `json:"address_ru,omitempty"`
	Latitude     float64     `json:"latitude,omitempty"`
	Longitude    float64     `json:"longitude,omitempty"`
	Image        null.String `json:"image,omitempty"`
	HasShipping  bool        `json:"has_shipping,omitempty"`
	ShopOwnerID  null.String `json:"shop_owner_id,omitempty"`
	ShopPhones   []string    `json:"phones,omitempty"`
	ParentShopID null.String `json:"-"`
	ParentShop   ParentShop  `json:"parent_shop,omitempty"`
}

type ParentShop struct {
	ID               string `json:"id,omitempty"`
	NameTM           string `json:"name_tm,omitempty"`
	NameRU           string `json:"name_ru,omitempty"`
	IsShoppingCenter string `json:"is_shopping_center,omitempty"`
}
