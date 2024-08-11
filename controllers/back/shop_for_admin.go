package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func GetAdminShops(c *gin.Context) {
	var shopQuery serializations.ShopQuery
	var shops []serializations.GetShops
	isDeleted := "NULL"
	selectedRows := "image"
	var queryShopOwner, search, searchStr, querySearch string

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request query - den maglumatlar bind edilyar
	if err := c.Bind(&shopQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&shopQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := shopQuery.Limit * (shopQuery.Page - 1)

	if shopQuery.Search != "" {
		incomingsSarch := slug.MakeLang(c.Query("search"), "en")
		search = strings.ReplaceAll(incomingsSarch, "-", " | ")
		searchStr = fmt.Sprintf("%%%s%%", search)
	}

	// request - den gelen deleted statusa gora pozulan ya-da pozulmadyk maglumatlar alynmaly
	if shopQuery.IsDeleted {
		isDeleted = "NOT NULL"
	}

	isShoppingCenter := shopQuery.IsShoppingCenter
	if isShoppingCenter {
		selectedRows = "latitude,longitude"
	}

	// request query - den status - a gora shop - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(
		`SELECT id,name_tm,name_ru,has_shipping,%s,created_status FROM shops WHERE deleted_at IS %v AND is_shopping_center=%v`,
		selectedRows, isDeleted, isShoppingCenter)

	if shopQuery.ShopOwnerID != "" {
		queryShopOwner = fmt.Sprintf(` AND shop_owner_id = '%v'`, shopQuery.ShopOwnerID)
	}

	if shopQuery.Search != "" {
		querySearch = fmt.Sprintf(` AND (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s')`, shopQuery.Lang, search, shopQuery.Lang, searchStr)
	}

	queryLimitOffset := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, shopQuery.Limit, offset)

	// database - den shop - lar alynyar
	var rowsShop pgx.Rows
	if len(shopQuery.CratedStatuses) != 0 {
		rowsShop, err = db.Query(context.Background(), rowQuery+queryShopOwner+querySearch+" AND created_status=ANY($1) "+queryLimitOffset, pq.Array(shopQuery.CratedStatuses))

	} else {
		rowsShop, err = db.Query(context.Background(), rowQuery+queryShopOwner+querySearch+queryLimitOffset)
	}
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	for rowsShop.Next() {
		var shop serializations.GetShops
		if isShoppingCenter {
			if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.HasShipping, &shop.Latitude, &shop.Longitude, &shop.CreatedStatus); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		} else {
			if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.HasShipping, &shop.Image, &shop.CreatedStatus); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}
