package models

type ProductColor struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"dimension_id" binding:"required"`
	ProductID string `json:"product_id,omitempty"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}
