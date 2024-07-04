package models

type ProductImage struct {
	ID             string `json:"id,omitempty"`
	Image          string `json:"image,omitempty"`
	ProductColorID string `json:"product_color_id,omitempty"`
	OrderNumber    int8   `json:"order_number,omitempty"`
}
