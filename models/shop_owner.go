package models

type ShopOwner struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty" binding:"required"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"required"`
	Password    string `json:"password,omitempty" binding:"required"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
	DeletedAt   string `json:"-"`
}
