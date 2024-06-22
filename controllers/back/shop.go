package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
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

	if err := models.ValidateShop(shop, true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda shops tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	var shop_id string
	if err = db.QueryRow(context.Background(), "INSERT INTO shops (name_tm,name_ru,address_tm,address_ru,latitude,longitude,image,has_delivery,shop_owner_id,slug_tm,slug_ru,order_number) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id", shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, shop.Image, shop.HasDelivery, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"), slug.MakeLang(shop.NameRU, "en"), shop.OrderNumber).Scan(&shop_id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// shop_phones tablisa maglumat gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO shop_phones (phone_number,shop_id) VALUES (unnest($1::varchar[]),$2)", pq.Array(shop.ShopPhones), shop_id)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// shop_categories tablisa maglumat gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO shop_categories (category_id,shop_id) VALUES (unnest($1::uuid[]),$2)", pq.Array(shop.Categories), shop_id)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
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

	if err := helpers.ValidateShopOwnerByToken(c, db, shop.ShopOwnerID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateShop(shop, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE shops SET name_tm=$1 , name_ru=$2 , address_tm=$3 , address_ru=$4 , latitude=$5 , longitude=$6 , image=$7 , has_delivery=$8 , shop_owner_id=$9 , slug_tm=$10 , slug_ru=$11 , order_number=$12 WHERE id=$13", shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, shop.Image, shop.HasDelivery, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"), slug.MakeLang(shop.NameRU, "en"), shop.OrderNumber, shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// shop - yn onki nomerleri pozulyar we taze telefon nomerler shop_phones tablisa gosulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shop_phones WHERE shop_id = $1", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	_, err = db.Exec(context.Background(), "INSERT INTO shop_phones (phone_number,shop_id) VALUES (unnest($1::varchar[]),$2)", pq.Array(shop.ShopPhones), shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// shop - a degisli onki kategoriyalar pozulyar we taze kategoriyalar gosulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shop_categories WHERE shop_id = $1", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	_, err = db.Exec(context.Background(), "INSERT INTO shop_categories (category_id,shop_id) VALUES (unnest($1::uuid[]),$2)", pq.Array(shop.Categories), shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
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
	var shop models.Shop
	if err := db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,address_tm,address_ru,latitude,longitude,image,has_delivery,shop_owner_id FROM shops WHERE id = $1 AND deleted_at IS NULL", shopID).Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shop.Image, &shop.HasDelivery, &shop.ShopOwnerID); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if shop.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	if err := helpers.ValidateShopOwnerByToken(c, db, shop.ShopOwnerID); err != nil {
		helpers.HandleError(c, 400, err.Error())
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

	// shop_id boyunca shop_category - lar cekilyar
	rowsShopCategory, err := db.Query(context.Background(), "SELECT category_id FROM shop_categories WHERE shop_id=$1 AND deleted_at IS NULL", shop.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShopCategory.Close()

	for rowsShopCategory.Next() {
		var categoryID string
		if err := rowsShopCategory.Scan(&categoryID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shop.Categories = append(shop.Categories, categoryID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})
}

func GetShops(c *gin.Context) {
	var shopQuery models.ShopQuery
	var shops []models.Shop
	isDeleted := "NULL"

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

	// request - den gelen deleted statusa gora pozulan ya-da pozulmadyk maglumatlar alynmaly
	if shopQuery.IsDeleted {
		isDeleted = "NOT NULL"
	}

	// request query - den status - a gora shop - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf("SELECT id,name_tm,name_ru,image FROM shops WHERE deleted_at IS %v ORDER BY created_at DESC LIMIT $1 OFFSET $2", isDeleted)
	if shopQuery.ShopOwnerID != "" {
		rows := strings.Split(rowQuery, " ORDER BY created_at DESC ")
		rowQuery = fmt.Sprintf("%v AND shop_owner_id = '%v' %v %v", rows[0], shopQuery.ShopOwnerID, "ORDER BY created_at DESC ", rows[1])
	}

	// database - den shop - lar alynyar
	rowsShop, err := db.Query(context.Background(), rowQuery, shopQuery.Limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShop.Close()

	for rowsShop.Next() {
		var shop models.Shop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.Image); err != nil {
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
