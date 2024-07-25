package models

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"strings"

	"github.com/gosimple/slug"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"
)

type Product struct {
	ID            string         `json:"id,omitempty"`
	NameTM        string         `json:"name_tm,omitempty" binding:"required"`
	NameRU        string         `json:"name_ru,omitempty" binding:"required"`
	Price         float32        `json:"price,omitempty" binding:"required"`
	OldPrice      null.Float     `json:"old_price,omitempty"`
	Code          string         `json:"code,omitempty"`
	SlugTM        string         `json:"slug_tm,omitempty"`
	SlugRU        string         `json:"slug_ru,omitempty"`
	BrendID       null.String    `json:"brend_id,omitempty"`
	ShopID        string         `json:"shop_id,omitempty" binding:"required"`
	Dimensions    []string       `json:"dimensions,omitempty"`
	Categories    []string       `json:"categories,omitempty" binding:"required"`
	ProductColors []ProductColor `json:"product_colors,omitempty" binding:"required"`
	Image         null.String    `json:"image,omitempty"`
	IsVisible     bool           `json:"is_visible,omitempty"`
}

func ValidateProduct(product Product, isCreateFunction bool) (productCode string, err error) {
	db, err := config.ConnDB()
	if err != nil {
		return "", err
	}
	defer db.Close()

	// harydyn bahasy we onki bahasy barlanyar
	if product.Price < 0 || product.OldPrice.Float64 < 0 {
		return "", errors.New("price or old_price cannot be less than 0")
	}
	if product.Price > float32(product.OldPrice.Float64) && product.OldPrice.Float64 != 0 {
		return "", errors.New("price cannot be less than old_price")
	}

	// eger haryda brend berilen bolsa onda sol brend hakykatdanam database - de barmy sol barlanyar
	if product.BrendID.String != "" {
		if err := helpers.ValidateRecordByID("brends", product.BrendID.String, "NULL", db); err != nil {
			return "", err
		}
	}

	// harydyn kategoriyalary barlanyar
	// hakykatdanam sol kategoriyalar barmy ?
	for _, v := range product.Categories {
		if err := helpers.ValidateRecordByID("categories", v, "NULL", db); err != nil {
			return "", err
		}
	}

	// harydyn razmerleri barlanyan
	// hakykatdanam sol razmerler database - de barmy ?
	for _, color := range product.ProductColors {
		for _, dimension := range color.Dimensions {
			if err := helpers.ValidateRecordByID("dimensions", dimension, "NULL", db); err != nil {
				return "", err
			}
		}
	}

	// haryt kot generate edilyar
	var categoryName, shopName string
	db.QueryRow(
		context.Background(),
		`SELECT name_tm FROM categories WHERE id=ANY($1) AND parent_category_id IS NULL AND deleted_at IS NULL`,
		pq.Array(product.Categories)).Scan(&categoryName)
	db.QueryRow(context.Background(), `SELECT name_tm FROM shops WHERE id=$1 AND deleted_at IS NULL`, product.ShopID).Scan(&shopName)
	code := strings.ToUpper(slug.MakeLang(shopName, "en")[:2]) + strings.ToUpper(slug.MakeLang(categoryName, "en")[:2]) + helpers.GenerateRandomCode()

	if !isCreateFunction {
		if product.ID == "" {
			return "", errors.New("id is required")
		}

		if err := helpers.ValidateRecordByID("products", product.ID, "NULL", db); err != nil {
			return "", err
		}
	}

	return code, nil
}
