package models

type ProductColor struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"dimension_id" binding:"required"`
	ProductID string `json:"product_id,omitempty"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
	DeletedAt string `json:"-"`
}

type ProductDimension struct {
	ID             string `json:"id,omitempty"`
	DimensionID    string `json:"dimension_id,omitempty"`
	ProductColorID string `json:"product_color_id,omitempty"`
	CreatedAt      string `json:"-"`
	UpdatedAt      string `json:"-"`
	DeletedAt      string `json:"-"`
}
