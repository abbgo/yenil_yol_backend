package models

type ProductImage struct {
	ID             string `json:"id,omitempty"`
	Image          string `json:"image,omitempty"`
	ProductColorID string `json:"product_color_id,omitempty"`
	CreatedAt      string `json:"-"`
	UpdatedAt      string `json:"-"`
	DeletedAt      string `json:"-"`
}
