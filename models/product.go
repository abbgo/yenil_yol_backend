package models

type Product struct {
	ID           string  `json:"id,omitempty"`
	NameTM       string  `json:"name_tm,omitempty" binding:"required"`
	NameRU       string  `json:"name_ru,omitempty" binding:"required"`
	Image        string  `json:"image,omitempty" binding:"required"`
	Price        float32 `json:"price,omitempty" binding:"required"`
	OldPrice     float32 `json:"old_price,omitempty"`
	Status       bool    `json:"status,omitempty"`
	ColorNameTM  string  `json:"color_name_tm,omitempty" binding:"required"`
	ColorNameRU  string  `json:"color_name_ru,omitempty" binding:"required"`
	GenderNameTM string  `json:"gender_name_tm,omitempty" binding:"required"`
	GenderNameRU string  `json:"gender_name_ru,omitempty" binding:"required"`
	Code         string  `json:"code,omitempty"`
	SlugTM       string  `json:"slug_tm,omitempty"`
	SlugRU       string  `json:"slug_ru,omitempty"`
	ShopID       string  `json:"shop_id,omitempty" binding:"required"`
	CategoryID   string  `json:"category_id,omitempty" binding:"required"`
	BrendID      string  `json:"brend_id,omitempty" binding:"required"`
	CreatedAt    string  `json:"-"`
	UpdatedAt    string  `json:"-"`
	DeletedAt    string  `json:"-"`
}
