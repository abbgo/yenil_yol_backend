package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

type ResponseShop struct {
	ID          string             `json:"id,omitempty"`
	NameTM      string             `json:"name_tm,omitempty"`
	NameRU      string             `json:"name_ru,omitempty"`
	AddressTM   string             `json:"address_tm,omitempty"`
	AddressRU   string             `json:"address_ru,omitempty"`
	Latitude    float64            `json:"latitude,omitempty"`
	Longitude   float64            `json:"longitude,omitempty"`
	Image       string             `json:"image,omitempty"`
	HasDelivery bool               `json:"has_delivery,omitempty"`
	ShopOwnerID string             `json:"shop_owner_id,omitempty"`
	SlugTM      string             `json:"slug_tm,omitempty"`
	SlugRU      string             `json:"slug_ru,omitempty"`
	ShopPhones  []models.ShopPhone `json:"shop_phones"`
}

func CreateShop(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var shop models.Shop
	if err := c.BindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//  request body - dan gelen telefon belgi barlanylyar
	for _, v := range shop.ShopPhones {
		if !helpers.ValidatePhoneNumber(v) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "invalid phone number",
			})
			return
		}
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if shop.Image == "" {
		image = nil
	} else {
		image = shop.Image
	}

	// eger maglumatlar dogry bolsa onda shops tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	var shop_id string
	if err = db.QueryRow(context.Background(), "INSERT INTO shops (name_tm,name_ru,address_tm,address_ru,latitude,longitude,image,has_delivery,shop_owner_id,slug_tm,slug_ru) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id", shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, image, shop.HasDelivery, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"), slug.MakeLang(shop.NameRU, "en")).Scan(&shop_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// shop_phones tablisa maglumat gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO shop_phones (phone_number,shop_id) VALUES (unnest($1::varchar[]),$2)", pq.Array(shop.ShopPhones), shop_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// shop - yn maglumatlary gosulandan sonra helper_images tablisa shop ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var shop models.Shop
	if err := c.BindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//  request body - dan gelen telefon belgi barlanylyar
	for _, v := range shop.ShopPhones {
		if !helpers.ValidatePhoneNumber(v) {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "invalid phone number",
			})
			return
		}
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	rowShop, err := db.Query(context.Background(), "SELECT id,image FROM shops WHERE id = $1 AND deleted_at IS NULL", shop.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer rowShop.Close()

	var oldShopImage, shopID string
	for rowShop.Next() {
		if err := rowShop.Scan(&shopID, &oldShopImage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if shopID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "shop not found",
		})
		return
	}

	// bu yerde shop - yn suraty ucin fileName atly uytgeyan ululyk doredilyar
	// eger request body - dan surat gelmese onda shop - yn suraty uytgedilmeyar diymek bolyar
	// sonun ucin shop - yn onki suratyny goyyas , eger request body - dan surat gelen bolsa
	// onda taze suraty kone surat bilen calysyas
	var fileName string
	if shop.Image == "" {
		fileName = oldShopImage
	} else {
		fileName = shop.Image

		// sonra helper_images tablisa shop ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", fileName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		// surat papkadan pozulyar
		if err := os.Remove(helpers.ServerPath + oldShopImage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE shops SET name_tm=$1 , name_ru=$2 , address_tm=$3 , address_ru=$4 , latitude=$5 , longitude=$6 , image=$7 , has_delivery=$8 , shop_owner_id=$9 , slug_tm=$10 , slug_ru=$11 , updated_at=$12 WHERE id=$13", shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, fileName, shop.HasDelivery, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"), slug.MakeLang(shop.NameRU, "en"), time.Now(), shop.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// shop - yn onli nomerleri pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shop_phones WHERE shop_id = $1", shop.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// shop_phones tablisa maglumat gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO shop_phones (phone_number,shop_id) VALUES (unnest($1::varchar[]),$2)", pq.Array(shop.ShopPhones), shop.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	shopID := c.Param("id")

	// database - den request parametr - den gelen id boyunca maglumat cekilyar
	rowShop, err := db.Query(context.Background(), "SELECT s.id,s.name_tm,s.name_ru,s.address_tm,s.address_ru,s.latitude,s.longitude,s.image,s.has_delivery,s.shop_owner_id,sp.id,sp.phone_number FROM shops s INNER JOIN shop_phones sp ON sp.shop_id = s.id WHERE s.id = $1 AND s.deleted_at IS NULL AND sp.deleted_at IS NULL", shopID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer rowShop.Close()

	var shop ResponseShop
	for rowShop.Next() {
		var shopPhone models.ShopPhone
		if err := rowShop.Scan(&shop.ID, &shop.NameTM, &shop.NameRU, &shop.AddressTM, &shop.AddressRU, &shop.Latitude, &shop.Longitude, &shop.Image, &shop.HasDelivery, &shop.ShopOwnerID, &shopPhone.ID, &shopPhone.PhoneNumber); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
				"error":   "yalnys bar",
			})
			return
		}
		shop.ShopPhones = append(shop.ShopPhones, shopPhone)
	}

	// eger databse sol maglumat yok bolsa error return edilyar
	if shop.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shop":   shop,
	})

}

func GetShops(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request parametr - den limit alynyar
	limitStr := c.Query("limit")
	if limitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "limit is required",
		})
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// request parametr - den page alynyar
	pageStr := c.Query("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "page is required",
		})
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// request parametr - den page alynyar
	shopOwnerID := c.DefaultQuery("shop_owner_id", "")

	// request query - den shop status alynyar
	// status -> shop pozulan ya-da pozulanmadygyny anlatyar
	// true bolsa pozulan
	// false bolsa pozulmadyk
	statusQuery := c.DefaultQuery("status", "false")
	status, err := strconv.ParseBool(statusQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// request query - den status - a gora shop - laryn sanyny almak ucin query yazylyar
	queryCount := fmt.Sprintf("SELECT COUNT(id) FROM shops WHERE deleted_at %v", "IS NULL")
	if status {
		queryCount = fmt.Sprintf("SELECT COUNT(id) FROM shops WHERE deleted_at %v", "IS NOT NULL")
	}

	if shopOwnerID != "" {
		queryCount = fmt.Sprintf("%v AND shop_owner_id = %v", queryCount, shopOwnerID)
	}

	// database - den shop - laryn sany alynyar
	var countOfShops uint
	if err = db.QueryRow(context.Background(), queryCount).Scan(&countOfShops); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "yalnys 1",
		})
		return
	}

	// request query - den status - a gora shop - lary almak ucin query yazylyar
	rowQuery := fmt.Sprintf("SELECT id,name_tm,image FROM shops WHERE deleted_at %v ORDER BY created_at DESC LIMIT $1 OFFSET $2", "IS NULL")
	if status {
		rowQuery = fmt.Sprintf("SELECT id,name_tm,image FROM shops WHERE deleted_at %v ORDER BY created_at DESC LIMIT $1 OFFSET $2", "IS NOT NULL")
	}

	if shopOwnerID != "" {
		// rowQuery = fmt.Sprintf("%v AND shop_owner_id = $%v", queryCount, shopOwnerID)
		rows := strings.Split(rowQuery, " ORDER BY created_at DESC ")
		rowQuery = fmt.Sprintf("%v AND shop_owner_id = %v %v %v", rows[0], shopOwnerID, "ORDER BY created_at DESC ", rows[1])
	}

	// database - den shop - lar alynyar
	rowsShop, err := db.Query(context.Background(), rowQuery, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "yalnys 2",
		})
		return
	}
	defer rowsShop.Close()

	var shops []ResponseShop
	for rowsShop.Next() {
		var shop ResponseShop
		if err := rowsShop.Scan(&shop.ID, &shop.NameTM, &shop.Image); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
		shops = append(shops, shop)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"shops":  shops,
		"total":  countOfShops,
	})

}

func DeleteShopByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM shops WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
		return
	}

	// hemme zat dogry bolsa shop we sol shop - yn we sol shop - a degisli shop_phones tablisalaryn deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_shop($1)", ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "yalnys bar",
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")

	// alynan id den bolan shop database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM shops WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// eger database sol id degisli shop yok bolsa error return edilyar
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
		return
	}

	// hemme zat dogry bolsa shop restore edilyar
	_, err = db.Exec(context.Background(), "CALL restore_shop($1)", ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var image string
	if err := db.QueryRow(context.Background(), "SELECT image FROM shops WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// eger database - de gelen id degisli shop yok bolsa error return edilyar
	if image == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "record not found",
		})
		return
	}

	// eger shop bar bolsa sonda shop - yn suraty papkadan pozulyar
	if err := os.Remove(helpers.ServerPath + image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// shop - yn suraty pozulandan sonra database - den shop pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shops WHERE id = $1", ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}
