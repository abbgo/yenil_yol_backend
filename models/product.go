package models

import (
	"errors"

	"gopkg.in/guregu/null.v4"
)

type Product struct {
	ID           string    `json:"id,omitempty"`
	NameTM       string    `json:"name_tm,omitempty" binding:"required"`
	NameRU       string    `json:"name_ru,omitempty" binding:"required"`
	Price        float32   `json:"price,omitempty" binding:"required"`
	OldPrice     float32   `json:"old_price,omitempty"`
	Status       null.Bool `json:"status,omitempty"`
	GenderNameTM string    `json:"gender_name_tm,omitempty" binding:"required"`
	GenderNameRU string    `json:"gender_name_ru,omitempty" binding:"required"`
	Code         string    `json:"code,omitempty"`
	SlugTM       string    `json:"slug_tm,omitempty"`
	SlugRU       string    `json:"slug_ru,omitempty"`
	BrendID      string    `json:"brend_id,omitempty" binding:"required"`
	Dimensions   []string  `json:"dimensions,omitempty"`
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
