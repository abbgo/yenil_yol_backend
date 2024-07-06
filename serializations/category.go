package serializations

import "gopkg.in/guregu/null.v4"

type CategoryQuery struct {
	ShopID string `form:"shop_id"`
}

type GetCategories struct {
	ID               string          `json:"id,omitempty"`
	NameTM           string          `json:"name_tm,omitempty" binding:"required"`
	NameRU           string          `json:"name_ru,omitempty" binding:"required"`
	ParentCategoryID null.String     `json:"parent_category_id,omitempty"`
	ChildCategories  []GetCategories `json:"child_categories"`
}
