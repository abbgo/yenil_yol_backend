package serializations

import (
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"

	"gopkg.in/guregu/null.v4"
)

type GetProducts struct {
	ID            string                `json:"id,omitempty"`
	NameTM        string                `json:"name_tm,omitempty" binding:"required"`
	NameRU        string                `json:"name_ru,omitempty" binding:"required"`
	Price         float32               `json:"price,omitempty" binding:"required"`
	OldPrice      null.Float            `json:"old_price,omitempty"`
	Code          string                `json:"code,omitempty"`
	SlugTM        string                `json:"slug_tm,omitempty"`
	SlugRU        string                `json:"slug_ru,omitempty"`
	BrendID       null.String           `json:"brend_id,omitempty"`
	ShopID        string                `json:"shop_id,omitempty" binding:"required"`
	Dimensions    []string              `json:"dimensions,omitempty"`
	Categories    []string              `json:"categories,omitempty" binding:"required"`
	ProductColors []models.ProductColor `json:"product_colors,omitempty" binding:"required"`
	Image         null.String           `json:"image"`
	Brend         BrendForProduct       `json:"brend,omitempty"`
	Shop          ShopForProduct        `json:"shop,omitempty"`
	IsVisible     bool                  `json:"is_visible,omitempty"`
}

type BrendForProduct struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
	Slug  string `json:"slug,omitempty"`
}

type ShopForProduct struct {
	ID     string `json:"id"`
	NameTM string `json:"name_tm"`
	NameRU string `json:"name_ru"`
}

type ProductQuery struct {
	helpers.StandartQuery
	ShopID     string   `form:"shop_id"`
	Categories []string `form:"categories"`
	ProductID  string   `form:"product_id"`
	Search     string   `form:"search"`
	Lang       string   `form:"lang"`
}
