package models

type Shop struct {
	ID        string  `json:"id,omitempty"`
	NameTM    string  `json:"name_tm,omitempty" binding:"required"`
	NameRU    string  `json:"name_ru,omitempty" binding:"required"`
	Address   string  `json:"address,omitempty" binding:"required"`
	Latitude  float32 `json:"latitude,omitempty" binding:"required"`
	Longitude float32 `json:"longitude,omitempty" binding:"required"`
	Image     string  `json:"image,omitempty"`
	Delivery  bool    `json:"delivery,omitempty" binding:"required"`
	CreatedAt string  `json:"-"`
	UpdatedAt string  `json:"-"`
	DeletedAt string  `json:"-"`
}
