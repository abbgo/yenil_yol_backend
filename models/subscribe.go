package models

type Subscribe struct {
	ID         string `json:"id,omitempty"`
	CustomerID string `json:"customer_id,omitempty" binding:"required"`
	ShopID     string `json:"shop_id,omitempty" binding:"required"`
}
