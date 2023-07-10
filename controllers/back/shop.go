package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

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
	_, err = db.Exec(context.Background(), "UPDATE shops SET name_tm=$1 , name_ru=$2 , address_tm=$3 , address_ru=$4 , latitude=$5 , longitude=$6 , image=$7 , has_delivery=$8 , shop_owner_id=$9 , slug_tm=$10 , slug_ru=$11 WHERE id=$12", shop.NameTM, shop.NameRU, shop.AddressTM, shop.AddressRU, shop.Latitude, shop.Longitude, fileName, shop.HasDelivery, shop.ShopOwnerID, slug.MakeLang(shop.NameTM, "en"), slug.MakeLang(shop.NameRU, "en"), shop.ID)
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
