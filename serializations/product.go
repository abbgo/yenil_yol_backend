package serializations

import (
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"

	"gopkg.in/guregu/null.v4"
)

type GetProductForBack struct {
	ID            string                `json:"id,omitempty"`
	NameTM        string                `json:"name_tm,omitempty"`
	NameRU        string                `json:"name_ru,omitempty"`
	Price         float32               `json:"price,omitempty"`
	OldPrice      null.Float            `json:"old_price,omitempty"`
	ShopID        string                `json:"shop_id,omitempty"`
	Categories    []GetCategories       `json:"categories,omitempty"`
	ProductColors []ProductColorForBack `json:"product_colors,omitempty"`
	IsVisible     bool                  `json:"is_visible,omitempty"`
	Brend         BrendForProduct       `json:"brend,omitempty"`
	BrendID       null.String           `json:"-"`
	Genders       []int                 `json:"genders"`
}

type GetProductsForBack struct {
	ID             string      `json:"id,omitempty"`
	NameTM         string      `json:"name_tm,omitempty"`
	NameRU         string      `json:"name_ru,omitempty"`
	Price          float32     `json:"price,omitempty"`
	OldPrice       null.Float  `json:"old_price,omitempty"`
	Image          null.String `json:"image"`
	IsVisible      bool        `json:"is_visible,omitempty"`
	CreatedStatus  int8        `json:"created_status"`
	RejectedReason null.String `json:"rejected_reason,omitempty"`
}

type GetProductsForFront struct {
	ID            string                `json:"id,omitempty"`
	NameTM        string                `json:"name_tm,omitempty"`
	NameRU        string                `json:"name_ru,omitempty"`
	Price         float32               `json:"price,omitempty"`
	OldPrice      null.Float            `json:"old_price,omitempty"`
	BrendID       null.String           `json:"brend_id,omitempty"`
	ShopID        string                `json:"shop_id,omitempty"`
	Dimensions    []string              `json:"dimensions,omitempty"`
	Categories    []string              `json:"categories,omitempty"`
	ProductColors []models.ProductColor `json:"product_colors,omitempty"`
	Image         null.String           `json:"image"`
	Brend         BrendForProduct       `json:"brend,omitempty"`
	Shop          ShopForProduct        `json:"shop,omitempty"`
	IsVisible     bool                  `json:"is_visible,omitempty"`
}

type GetProductsForAdminProduct struct {
	ID            string                 `json:"id"`
	NameTM        string                 `json:"name_tm"`
	NameRU        string                 `json:"name_ru"`
	Price         float32                `json:"price"`
	OldPrice      null.Float             `json:"old_price,omitempty"`
	BrendID       null.String            `json:"-"`
	ShopID        string                 `json:"-"`
	Brend         BrendForProduct        `json:"brend,omitempty"`
	Shop          ShopForProduct         `json:"shop"`
	IsVisible     bool                   `json:"is_visible"`
	Categories    []CategoryForProduct   `json:"categories"`
	ProductColors []ProductColorForAdmin `json:"product_colors"`
	Genders       []int                  `json:"genders"`
}

type BrendForProduct struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
	Slug  string `json:"slug,omitempty"`
}

type ShopForProduct struct {
	ID      string `json:"id"`
	NameTM  string `json:"name_tm"`
	NameRU  string `json:"name_ru"`
	IsBrand bool   `json:"is_brand"`
}

type ProductQuery struct {
	helpers.StandartQuery
	ShopID         string   `form:"shop_id"`
	Categories     []string `form:"categories"`
	ProductID      string   `form:"product_id"`
	Search         string   `form:"search"`
	Lang           string   `form:"lang"`
	Sort           string   `form:"sort"`
	MinPrice       string   `form:"min_price"`
	MaxPrice       string   `form:"max_price"`
	CratedStatuses []string `form:"crated_statuses"`
	Genders        []string `form:"genders"`
}

type ProductCountQuery struct {
	IsDeleted bool   `form:"is_deleted"`
	ShopID    string `form:"shop_id"`
	Search    string `form:"search"`
	Lang      string `form:"lang"`
}
