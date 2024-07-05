package serializations

import "github/abbgo/yenil_yol/backend/helpers"

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
