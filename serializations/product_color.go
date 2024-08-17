package serializations

import (
	"github/abbgo/yenil_yol/backend/models"

	"gopkg.in/guregu/null.v4"
)

type ProductColorForBack struct {
	ID          string                `json:"id,omitempty"`
	Name        null.String           `json:"name"`
	Images      []models.ProductImage `json:"images,omitempty"`
	Dimensions  []models.Dimension    `json:"dimensions,omitempty"`
	OrderNumber int8                  `json:"order_number,omitempty"`
}

type ProductColorForAdmin struct {
	ID         string      `json:"-"`
	Name       null.String `json:"name,omitempty"`
	Images     []string    `json:"images"`
	Dimensions []string    `json:"dimensions"`
}
