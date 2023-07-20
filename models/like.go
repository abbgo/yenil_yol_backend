package models

type Like struct {
	ID         string `json:"id,omitempty"`
	CustomerID string `json:"customer_id,omitempty" binding:"required"`
	ProductID  string `json:"product_id,omitempty" binding:"required"`
	CreatedAt  string `json:"-"`
	UpdatedAt  string `json:"-"`
	DeletedAt  string `json:"-"`
}
