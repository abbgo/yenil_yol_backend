package models

type ProductColor struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"dimension_id" binding:"required"`
	ProductID string `json:"product_id,omitempty"`
}

type ProductDimension struct {
	ID             string `json:"id,omitempty"`
	DimensionID    string `json:"dimension_id,omitempty"`
	ProductColorID string `json:"product_color_id,omitempty"`
}
