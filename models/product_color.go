package models

import "gopkg.in/guregu/null.v4"

type ProductColor struct {
	ID         string         `json:"id,omitempty"`
	Name       null.String    `json:"name" binding:"required"`
	ProductID  string         `json:"product_id,omitempty"`
	Images     []ProductImage `json:"images,omitempty" binding:"required"`
	Dimensions []string       `json:"dimensions,omitempty" binding:"required"`
}

type ProductDimension struct {
	ID             string `json:"id,omitempty"`
	DimensionID    string `json:"dimension_id,omitempty"`
	ProductColorID string `json:"product_color_id,omitempty"`
}
