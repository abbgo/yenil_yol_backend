package models

import (
	"errors"

	"gopkg.in/guregu/null.v4"
)

type Product struct {
	ID            string         `json:"id,omitempty"`
	NameTM        string         `json:"name_tm,omitempty" binding:"required"`
	NameRU        string         `json:"name_ru,omitempty" binding:"required"`
	Price         float32        `json:"price,omitempty" binding:"required"`
	OldPrice      float32        `json:"old_price,omitempty"`
	Code          string         `json:"code,omitempty"`
	SlugTM        string         `json:"slug_tm,omitempty"`
	SlugRU        string         `json:"slug_ru,omitempty"`
	BrendID       null.String    `json:"brend_id,omitempty"`
	Dimensions    []string       `json:"dimensions,omitempty"`
	Categories    []string       `json:"categories,omitempty" binding:"required"`
	ProductColors []ProductColor `json:"product_colors,omitempty" binding:"required"`
}

func ValidateProduct(product Product) error {
	if product.Price < 0 || product.OldPrice < 0 {
		return errors.New("price or old_price cannot be less than 0")
	}

	if product.Price > product.OldPrice && product.OldPrice != 0 {
		return errors.New("price cannot be less than old_price")
	}

	return nil
}
