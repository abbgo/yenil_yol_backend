package controllers

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterShopOwner(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request - den gelen maglumatlar alynyar
	var shopOwner models.ShopOwner
	if err := c.BindJSON(&shopOwner); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// gelen maglumatlar barlanylyar
	if err := models.ValidateRegisterShopOwner(shopOwner.PhoneNumber, "register"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// parol hashlenyan
	hashPassword, err := helpers.HashPassword(shopOwner.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// hemme zat yerbe yer bolsa maglumatlar shop_owners tablisa gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO shop_owners (name,phone_number,password) VALUES ($1,$2,$3)", shopOwner.Name, shopOwner.PhoneNumber, hashPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "1",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"phone_number": shopOwner.PhoneNumber,
		"password":     shopOwner.Password,
	})

}

func LoginShopOwner(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request - den maglumatlar alynyar
	var shopOwner models.ShopOwnerLogin
	if err := c.BindJSON(&shopOwner); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := models.ValidateRegisterShopOwner(shopOwner.PhoneNumber, "login"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// database - den telefon belgisi request - den gelen telefon belga den bolan maglumat cekilyar
	var id, password string
	row, err := db.Query(context.Background(), "SELECT id,password FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", shopOwner.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer row.Close()

	for row.Next() {
		if err := row.Scan(&id, &password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
	}

	// eger request - den gelen telefon belgili shop_owner database - de yok bolsa onda error response edilyar
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "this show_owner does not exist",
		})
		return
	}

	// eger shop_owner bar bolsa onda paroly dogry yazypdyrmy yazmandyrmy sol barlanyar
	credentialError := helpers.CheckPassword(shopOwner.Password, password)
	if credentialError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid credentials",
		})
		return
	}

	// maglumatlar dogry bolsa auth ucin access_toke bilen resfresh_token generate edilyar
	accessTokenString, refreshTokenString, err := helpers.GenerateAccessTokenForAdmin(shopOwner.PhoneNumber, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// front tarapa ugratmak ucin shop_owner - in id - si boyunca maglumatlary get edilyar
	adm, err := GetShopOwnerByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
		"admin":         adm,
	})

}

func GetShopOwnerByID(id string) (models.ShopOwner, error) {
	db, err := config.ConnDB()
	if err != nil {
		return models.ShopOwner{}, err
	}
	defer db.Close()

	// parametrler edilip berilen id - boyunca database - den shop_owner - in maglumatlary cekilyar
	var shopOwner models.ShopOwner
	rowShopOwner, err := db.Query(context.Background(), "SELECT name,phone_number FROM shop_owners WHERE deleted_at IS NULL AND id = $1", id)
	if err != nil {
		return models.ShopOwner{}, err
	}
	defer rowShopOwner.Close()

	for rowShopOwner.Next() {
		if err := rowShopOwner.Scan(&shopOwner.Name, &shopOwner.PhoneNumber); err != nil {
			return models.ShopOwner{}, err
		}
	}

	// eger parametrler edilip berilen id boyunca database - de maglumat yok bolsa error return edilyar
	if shopOwner.PhoneNumber == "" {
		return models.ShopOwner{}, errors.New("shop_owner not found")
	}

	// hemme zat dogry bolsa shop_owner - in maglumatlary return edilyar
	return shopOwner, nil

}

func UpdateShopOwner(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - den shop_owner - in maglumatlary alynyar
	var shopOwner models.ShopOwnerUpdate
	if err := c.BindJSON(&shopOwner); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// database - de request body - den gelen id bilen gabat gelyan shop_owner barmy ya-da yokmy sol barlanyar
	// eger yok bolsa onda error return edilyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM shop_owners WHERE id = $1 AND deleted_at IS NULL", shopOwner.ID).Scan(&id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "shop_owner not found",
		})
		return
	}

	// eger shop_owner database - de bar bolsa onda onun maglumatlary request body - dan gelen maglumatlar bilen update edilyar
	_, err = db.Exec(context.Background(), "UPDATE shop_owners SET name = $1 , phone_number = $2 WHERE id = $3", shopOwner.Name, shopOwner.PhoneNumber, id)
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

func GetShopOwners(c *gin.Context) {

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
	limitStr := c.Param("limit")
	if limitStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "limit is required",
		})
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// // request parametr - den page alynyar
	pageStr := c.Param("page")
	if pageStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "page is required",
		})
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// database - den shop_owner - lerin sany alynyar
	countOfShopOwners := 0
	if err := db.QueryRow(context.Background(), "SELECT COUNT(id) FROM shop_owners WHERE deleted_at IS NULL").Scan(&countOfShopOwners); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// databae - den request - den gelen limit we page boyunca limitlap shop_owner - ler alynyar
	var shopOwners []models.ShopOwner
	rowsShopOwner, err := db.Query(context.Background(), "SELECT name,phone_number FROM shop_owners WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer rowsShopOwner.Close()

	for rowsShopOwner.Next() {
		var shopOwner models.ShopOwner
		if err := rowsShopOwner.Scan(&shopOwner.Name, &shopOwner.PhoneNumber); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
		shopOwners = append(shopOwners, shopOwner)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"shop_owners": shopOwners,
		"total":       countOfShopOwners,
	})

}

func GetShopOwnerD(c *gin.Context) {

	shopOwnerID, hasID := c.Get("id")
	if !hasID {
		c.JSON(http.StatusBadRequest, "shopOwnerID is required")
		return
	}

	var ok bool
	shopOwner_id, ok := shopOwnerID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, "shopOwnerID must be uint")
	}

	adm, err := GetShopOwnerByID(shopOwner_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shop_owner": adm,
	})

}
