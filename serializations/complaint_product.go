package serializations

import "github/abbgo/yenil_yol/backend/helpers"

type ComplaintProduct struct {
	ID             string `json:"id,omitempty"`
	NameTM         string `json:"name_tm,omitempty"`
	NameRU         string `json:"name_ru,omitempty"`
	Image          string `json:"image,omitempty"`
	ComplaintCount int    `json:"complaint_count,omitempty"`
}

type ComplaintProductQuery struct {
	helpers.StandartQuery
	ShopOwnerID string `form:"shop_owner_id"`
}
