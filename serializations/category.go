package serializations

import (
	"github/abbgo/yenil_yol/backend/helpers"

	"gopkg.in/guregu/null.v4"
)

type CategoryQuery struct {
	helpers.StandartQuery
	ShopID string `form:"shop_id"`
	Search string `form:"search"`
	Lang   string `form:"lang"`
}

type GetCategories struct {
	ID               string          `json:"id,omitempty"`
	NameTM           string          `json:"name_tm,omitempty"`
	NameRU           string          `json:"name_ru,omitempty"`
	ParentCategoryID null.String     `json:"parent_category_id,omitempty"`
	ChildCategories  []GetCategories `json:"child_categories,omitempty"`
}
type GetCategoriesForAdmin struct {
	ID               string                  `json:"id,omitempty"`
	NameTM           string                  `json:"name_tm,omitempty"`
	NameRU           string                  `json:"name_ru,omitempty"`
	Image            null.String             `json:"image,omitempty"`
	ParentCategoryID null.String             `json:"parent_category_id,omitempty"`
	ChildCategories  []GetCategoriesForAdmin `json:"child_categories,omitempty"`
}

type CategoryForProduct struct {
	ID     string `json:"id"`
	NameTM string `json:"name_tm"`
	NameRU string `json:"name_ru"`
}
