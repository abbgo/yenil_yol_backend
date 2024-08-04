package models

type ComplaintProduct struct {
	ID          string `json:"id,omitempty"`
	ComplaintID string `json:"complaint_id,omitempty" binding:"required"`
	ProductID   string `json:"product_id,omitempty" binding:"required"`
}
