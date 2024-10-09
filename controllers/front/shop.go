package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

func GetShopsForMap(c *gin.Context) {
	var requestQuery serializations.ShopForMapQuery
	var gendersQuery, joinProductsQuery string

	// request query - den maglumatlar bind edilyar
	if err := c.Bind(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if len(requestQuery.Genders) != 0 {
		joinProductsQuery = ` INNER JOIN products p ON s.id = p.shop_id `

		gender, err := strconv.ParseInt(requestQuery.Genders[0], 10, 8)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		gendersQuery += fmt.Sprintf(` AND (%d = ANY(p.genders) `, gender)

		if len(requestQuery.Genders) > 1 {
			genders := requestQuery.Genders[1:]
			lenghtGenders := len(genders)
			for i := 0; i < lenghtGenders; i++ {
				gender, err := strconv.ParseInt(genders[i], 10, 8)
				if err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}
				if genders[i] == genders[lenghtGenders-1] {
					gendersQuery += fmt.Sprintf(` OR %d = ANY(p.genders)) `, gender)
				} else {
					gendersQuery += fmt.Sprintf(` OR %d = ANY(p.genders) `, gender)
				}
			}
		} else {
			gendersQuery += `)`
		}
	}

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	rowsShopQuery := fmt.Sprintf(`
							SELECT DISTINCT ON (s.id) s.id,s.name_tm,s.name_ru,s.latitude,s.longitude,s.is_shopping_center FROM shops s  
							%s
							WHERE 6371 * acos(
										cos(radians(%f)) * cos(radians(s.latitude)) *
										cos(radians(s.longitude) - radians(%f)) +
										sin(radians(%f)) * sin(radians(s.latitude))
									) <= %d AND s.deleted_at IS NULL AND s.parent_shop_id IS NULL AND 
									 s.created_status=%d AND s.at_home=false %s;
						`, joinProductsQuery, requestQuery.Latitude, requestQuery.Longitude, requestQuery.Latitude,
		requestQuery.Kilometer, helpers.CreatedStatuses["success"], gendersQuery)

	rowsShop, err := db.Query(context.Background(), rowsShopQuery)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()
	var shops []models.Shop
	for rowsShop.Next() {
		var shop models.Shop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude, &shop.IsShoppingCenter); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}

func GetShops(c *gin.Context) {
	var requestQuery serializations.ShopQuery
	var search, searchStr, querySearch, parentShopIDQuery string
	queryRandom := ` ORDER BY created_at DESC`

	// request query - den maglumatlar bind edilyar
	if err := c.Bind(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := requestQuery.Limit * (requestQuery.Page - 1)

	if requestQuery.Search != "" {
		incomingsSarch := slug.MakeLang(c.Query("search"), "en")
		search = strings.ReplaceAll(incomingsSarch, "-", " | ")
		searchStr = fmt.Sprintf("%%%s%%", search)
	}

	if requestQuery.ParentShopID != "" {
		parentShopIDQuery = fmt.Sprintf(` AND parent_shop_id='%s'`, requestQuery.ParentShopID)
	}

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// database - den shop - lar alynyar
	queryDefault := fmt.Sprintf(
		`SELECT id,name_tm,name_ru,latitude,longitude,resized_image,address_tm,address_ru,parent_shop_id,is_shopping_center,at_home,is_brand FROM shops WHERE 
		deleted_at IS NULL AND created_status=%d AND (is_shopping_center=false OR is_shopping_center=%v)`, helpers.CreatedStatuses["success"],
		requestQuery.IsShoppingCenter)

	if requestQuery.IsRandom {
		queryRandom = ` ORDER BY RANDOM()`
	}
	if requestQuery.Search != "" {
		querySearch = fmt.Sprintf(` AND (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s')`, requestQuery.Lang, search, requestQuery.Lang, searchStr)
	}
	queryLimitOffset := fmt.Sprintf(` LIMIT %v OFFSET %v`, requestQuery.Limit, offset)

	rowsShop, err := db.Query(context.Background(), queryDefault+parentShopIDQuery+querySearch+queryRandom+queryLimitOffset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	var shops []serializations.GetShop
	for rowsShop.Next() {
		var shop serializations.GetShop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Latitude, &shop.Longitude, &shop.Image,
			&shop.AddressTM, &shop.AddressRU, &shop.ParentShopID, &shop.IsShoppingCenter, &shop.AtHome, &shop.IsBrand); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// eger shop - yn parent shop - y bar bolsa onda sol alynyar
		if shop.ParentShopID.String != "" {
			var parentShop serializations.ParentShop
			if err := db.QueryRow(
				context.Background(),
				`SELECT id,name_tm,name_ru,is_shopping_center FROM shops WHERE id=$1 AND created_status=$2`,
				shop.ParentShopID.String, helpers.CreatedStatuses["success"],
			).
				Scan(&parentShop.ID, &parentShop.NameTM, &parentShop.NameRU, &parentShop.IsShoppingCenter); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			shop.ParentShop = &parentShop
		}

		rowsShopPhones, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id = $1 AND deleted_at IS NULL", shop.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsShopPhones.Close()

		for rowsShopPhones.Next() {
			var phoneNumber string
			if err := rowsShopPhones.Scan(&phoneNumber); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			shop.ShopPhones = append(shop.ShopPhones, phoneNumber)
		}
		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}

func GetShopByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	shopID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	var shop models.Shop
	db.QueryRow(
		context.Background(), `SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,resized_image,has_shipping,is_brand FROM shops 
		WHERE id = $1 AND created_status=$2 AND deleted_at IS NULL`,
		shopID, helpers.CreatedStatuses["success"]).
		Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shop.Image, &shop.HasShipping, &shop.IsBrand)

	// eger databse sol maglumat yok bolsa error return edilyar
	if shop.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// shop - a degisli telefon belgiler alynyar
	rowsShopPhone, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id=$1 AND deleted_at IS NULL", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShopPhone.Close()
	for rowsShopPhone.Next() {
		var phoneNumber string
		if err := rowsShopPhone.Scan(&phoneNumber); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shop.ShopPhones = append(shop.ShopPhones, phoneNumber)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})
}

func GetShopByIDs(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id - ler alynyar
	shopIDs := c.QueryArray("ids")

	// database - den request parametr - den gelen id - ler boyunca maglumat cekilyar
	var shops []serializations.GetShop
	rows, err := db.Query(context.Background(),
		`
			SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,resized_image,parent_shop_id,at_home,is_brand FROM shops 
			WHERE id = ANY($1) AND created_status=$2 AND deleted_at IS NULL AND is_shopping_center=false
		`,
		pq.Array(shopIDs), helpers.CreatedStatuses["success"])
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var shop serializations.GetShop
		if err := rows.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU,
			&shop.Latitude, &shop.Longitude, &shop.Image, &shop.ParentShopID, &shop.AtHome, &shop.IsBrand); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if shop.ParentShopID.String != "" {
			var parentShop serializations.ParentShop
			db.QueryRow(context.Background(), `SELECT id,name_tm,name_ru,is_shopping_center FROM shops WHERE id=$1`, shop.ParentShopID.String).
				Scan(&parentShop.ID, &parentShop.NameTM, &parentShop.NameRU, &parentShop.IsShoppingCenter)

			shop.ParentShop = &parentShop
		}

		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
	})
}
