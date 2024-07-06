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
}

type ShopForMapQuery struct {
	Latitude  float64 `form:"latitude" validate:"required"`
	Longitude float64 `form:"longitude" validate:"required"`
	Kilometer int8    `form:"kilometer"`
}

type GetShops struct {
	ID          string      `json:"id,omitempty"`
	NameTM      string      `json:"name_tm,omitempty" binding:"required"`
	NameRU      string      `json:"name_ru,omitempty" binding:"required"`
	Latitude    float64     `json:"latitude,omitempty" binding:"required"`
	Longitude   float64     `json:"longitude,omitempty" binding:"required"`
	Image       null.String `json:"image,omitempty"`
	HasShipping bool        `json:"has_shipping"`
}

type GetShop struct {
	ID           string      `json:"id,omitempty"`
	NameTM       string      `json:"name_tm,omitempty" binding:"required"`
	NameRU       string      `json:"name_ru,omitempty" binding:"required"`
	AddressTM    string      `json:"address_tm,omitempty" binding:"required"`
	AddressRU    string      `json:"address_ru,omitempty" binding:"required"`
	Latitude     float64     `json:"latitude,omitempty" binding:"required"`
	Longitude    float64     `json:"longitude,omitempty" binding:"required"`
	Image        null.String `json:"image,omitempty"`
	HasShipping  bool        `json:"has_shipping,omitempty"`
	ShopOwnerID  null.String `json:"shop_owner_id,omitempty"`
	ShopPhones   []string    `json:"phones,omitempty"`
	ParentShopID null.String `json:"-"`
	ParentShop   ParentShop  `json:"parent_shop,omitempty"`
}

type ParentShop struct {
	NameTM string `json:"name_tm"`
	NameRU string `json:"name_ru"`
}
