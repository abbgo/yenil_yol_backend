package serializations

import "github/abbgo/yenil_yol/backend/helpers"

type CategoryQuery struct {
	helpers.StandartQuery
	ShopID string `form:"shop_id"`
}
