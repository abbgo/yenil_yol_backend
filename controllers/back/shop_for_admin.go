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
	var shops []serializations.GetShop
	isDeleted := "NULL"
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

	// request query - den status - a gora shop - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf(`SELECT id,image,name_tm,name_ru,address_tm,address_ru,latitude,longitude,has_shipping,shop_owner_id,parent_shop_id FROM shops 
	WHERE deleted_at IS %v AND is_shopping_center=false`, isDeleted)

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
		var shop serializations.GetShop
		if err := rowsShop.Scan(
			&shop.ID, &shop.Image, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude,
			&shop.HasShipping, &shop.ShopOwnerID, &shop.ParentShopID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// shop alynanadan son shop_id boyunca shop_phone - lar cekilyar
		rowsPhoneNumber, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id=$1 AND deleted_at IS NULL", shop.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsPhoneNumber.Close()

		for rowsPhoneNumber.Next() {
			var phoneNumber string
			if err := rowsPhoneNumber.Scan(&phoneNumber); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			shop.ShopPhones = append(shop.ShopPhones, phoneNumber)
		}

		if shop.ParentShopID.String != "" {
			var parentShop serializations.ParentShop
			if err := db.QueryRow(context.Background(), `SELECT id,name_tm,name_ru FROM shops WHERE id=$1`, shop.ParentShopID.String).
				Scan(&parentShop.ID, &parentShop.NameTM, &parentShop.NameRU); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			shop.ParentShop = parentShop
		}

		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}
