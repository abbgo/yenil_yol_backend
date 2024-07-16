package serializations

import "gopkg.in/guregu/null.v4"

type CategoryQuery struct {
	ShopID string `form:"shop_id"`
	Search string `form:"search"`
}

type GetCategories struct {
	ID               string          `json:"id,omitempty"`
	NameTM           string          `json:"name_tm,omitempty"`
	NameRU           string          `json:"name_ru,omitempty"`
	ParentCategoryID null.String     `json:"parent_category_id,omitempty"`
	ChildCategories  []GetCategories `json:"child_categories,omitempty"`
}
