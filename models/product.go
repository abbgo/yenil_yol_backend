package models

import (
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"

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
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// harydyn bahasy we onki bahasy barlanyar
	if product.Price < 0 || product.OldPrice < 0 {
		return errors.New("price or old_price cannot be less than 0")
	}
	if product.Price > product.OldPrice && product.OldPrice != 0 {
		return errors.New("price cannot be less than old_price")
	}

	// eger haryda brend berilen bolsa onda sol brend hakykatdanam database - de barmy sol barlanyar
	if product.BrendID.String != "" {
		if err := helpers.ValidateRecordByID("brends", product.BrendID.String, "NULL", db); err != nil {
			return err
		}
	}

	// harydyn kategoriyalary barlanyar
	// hakykatdanam sol kategoriyalar barmy ?
	for _, v := range product.Categories {
		if err := helpers.ValidateRecordByID("categories", v, "NULL", db); err != nil {
			return err
		}
	}

	// harydyn razmerleri barlanyan
	// hakykatdanam sol razmerler database - de barmy ?
	for _, color := range product.ProductColors {
		for _, dimension := range color.Dimensions {
			if err := helpers.ValidateRecordByID("dimensions", dimension, "NULL", db); err != nil {
				return err
			}
		}
	}

	return nil
}
