package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"github/abbgo/yenil_yol/backend/serializations"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

func CreateShop(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var shop models.Shop
	if err := c.BindJSON(&shop); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCreateShop(shop); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	var resizedImage interface{}
	if shop.Image.String == "" {
		resizedImage = nil
	} else {
		resizedImage = "assets/" + shop.Image.String
	}

	// eger maglumatlar dogry bolsa onda shops tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	var shop_id string
	if err = db.QueryRow(context.Background(),
		`INSERT INTO shops 
		(name_tm,name_ru,address_tm,address_ru,latitude,longitude,image,has_shipping,shop_owner_id,slug_tm,slug_ru,order_number,parent_shop_id,is_shopping_center,resized_image,at_home) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id`,
		shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, shop.Image, shop.HasShipping, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"),
		slug.MakeLang(shop.NameRU, "en"), shop.OrderNumber, shop.ParentShopID, shop.IsShoppingCenter, resizedImage, shop.AtHome).
		Scan(&shop_id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger shop sowda merkezi dal bolsa shop_phones tablisa maglumat gosulyar
	if len(shop.ShopPhones) != 0 {
		_, err = db.Exec(context.Background(), "INSERT INTO shop_phones (phone_number,shop_id) VALUES (unnest($1::varchar[]),$2)", pq.Array(shop.ShopPhones), shop_id)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// shop - yn maglumatlary gosulandan sonra helper_images tablisa shop ucin gosulan surat pozulyar
	if shop.Image.String != "" {
		if err := DeleteImageFromDB(shop.Image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})
}

func UpdateShopByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var shop models.Shop
	if err := c.BindJSON(&shop); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateUpdateShop(shop); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	var resizedImage interface{}
	if shop.Image.String == "" {
		resizedImage = nil
	} else {
		resizedImage = "assets/" + shop.Image.String
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(),
		`UPDATE shops SET name_tm=$1 , name_ru=$2 , address_tm=$3 , address_ru=$4 , latitude=$5 , longitude=$6 , image=$7 , has_shipping=$8 , shop_owner_id=$9 , slug_tm=$10 , slug_ru=$11 , 
		order_number=$12 , parent_shop_id=$13 , is_shopping_center=$14 , resized_image=$15 , created_status=0 , rejected_reason=NULL , at_home=$16 WHERE id=$17`,
		shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, shop.Image, shop.HasShipping, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"),
		slug.MakeLang(shop.NameRU, "en"), shop.OrderNumber, shop.ParentShopID, shop.IsShoppingCenter, resizedImage, shop.AtHome, shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger shop sowda merkezi dal bolsa shop - yn onki nomerleri pozulyar we taze telefon nomerler shop_phones tablisa gosulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shop_phones WHERE shop_id = $1", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	if len(shop.ShopPhones) != 0 {
		_, err = db.Exec(context.Background(), "INSERT INTO shop_phones (phone_number,shop_id) VALUES (unnest($1::varchar[]),$2)", pq.Array(shop.ShopPhones), shop.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	if shop.Image.String != "" {
		if err := DeleteImageFromDB(shop.Image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
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

	// database - den request parametr - den gelen id boyunca shop - yn maglumatlary cekilyar
	var shop serializations.GetShop
	if err := db.QueryRow(
		context.Background(),
		`SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,image,has_shipping,shop_owner_id,parent_shop_id,at_home 
		FROM shops WHERE id = $1 AND deleted_at IS NULL`,
		shopID).
		Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU,
			&shop.Latitude, &shop.Longitude, &shop.Image,
			&shop.HasShipping, &shop.ShopOwnerID, &shop.ParentShopID, &shop.AtHome); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if shop.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// shop alynanadan son shop_id boyunca shop_phone - lar cekilyar
	rowsShopImage, err := db.Query(context.Background(), "SELECT phone_number FROM shop_phones WHERE shop_id=$1 AND deleted_at IS NULL", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShopImage.Close()

	for rowsShopImage.Next() {
		var phoneNumber string
		if err := rowsShopImage.Scan(&phoneNumber); err != nil {
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
		shop.ParentShop = &parentShop
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})
}

func GetShops(c *gin.Context) {
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
		`SELECT id,name_tm,name_ru,has_shipping,%s,created_status,rejected_reason FROM shops WHERE deleted_at IS %v AND is_shopping_center=%v`,
		selectedRows, isDeleted, isShoppingCenter)

	if shopQuery.ShopOwnerID != "" {
		queryShopOwner = fmt.Sprintf(` AND shop_owner_id = '%v'`, shopQuery.ShopOwnerID)
	}

	if shopQuery.Search != "" {
		querySearch = fmt.Sprintf(` AND (to_tsvector(slug_%s) @@ to_tsquery('%s') OR slug_%s LIKE '%s')`, shopQuery.Lang, search, shopQuery.Lang, searchStr)
	}

	queryLimitOffset := fmt.Sprintf(` ORDER BY created_at DESC LIMIT %v OFFSET %v`, shopQuery.Limit, offset)

	// database - den shop - lar alynyar
	rowsShop, err := db.Query(context.Background(), rowQuery+queryShopOwner+querySearch+queryLimitOffset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	for rowsShop.Next() {
		var shop serializations.GetShops
		if isShoppingCenter {
			if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.HasShipping, &shop.Latitude, &shop.Longitude, &shop.CreatedStatus, &shop.RejectedReason); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		} else {
			if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.HasShipping, &shop.Image, &shop.CreatedStatus, &shop.RejectedReason); err != nil {
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

func DeleteShopByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("shops", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa shop we sol shop - yn we sol shop - a degisli shop_phones tablisalaryn deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_shop($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreShopByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("shops", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa shop restore edilyar
	_, err = db.Exec(context.Background(), "CALL restore_shop($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyShopByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	var image sql.NullString
	db.QueryRow(context.Background(), "SELECT id,image FROM shops WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id, &image)

	// eger database - de gelen id degisli shop yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger shop - yn suraty bar bolsa onda ol papkadan pozulyar
	if image.String != "" {
		if err := os.Remove(helpers.ServerPath + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := os.Remove(helpers.ServerPath + "assets/" + image.String); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// shop - yn suraty pozulandan sonra shop ve sonun bilen baglanysykly maglumatlar pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shops WHERE id=$1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}
