package models

type Shop struct {
	ID          string  `json:"id,omitempty"`
	NameTM      string  `json:"name_tm,omitempty" binding:"required"`
	NameRU      string  `json:"name_ru,omitempty" binding:"required"`
	Address     string  `json:"address,omitempty" binding:"required"`
	Latitude    float32 `json:"latitude,omitempty" binding:"required"`
	Longitude   float32 `json:"longitude,omitempty" binding:"required"`
	Image       string  `json:"image,omitempty"`
	HasDelivery bool    `json:"has_delivery,omitempty" binding:"required"`
	ShopOwnerID string  `json:"shop_owner_id,omitempty" binding:"required"`
	SlugTM      string  `json:"slug_tm,omitempty"`
	SlugRU      string  `json:"slug_ru,omitempty"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

type ShopPhone struct {
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	ShopID      string `json:"shop_id,omitempty" binding:"required"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}
