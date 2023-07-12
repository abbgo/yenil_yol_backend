package controllers

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAdmin(c *gin.Context) {

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
	var admin models.Admin
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// gelen telefon belgi barlanylyar
	if !helpers.ValidatePhoneNumber(admin.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid phone number",
		})
		return
	}

	// parol hashlenyan
	hashPassword, err := helpers.HashPassword(admin.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// hemme zat yerbe yer bolsa maglumatlar shop_owners tablisa gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO admins (full_name,phone_number,password,is_super_admin) VALUES ($1,$2,$3,$4)", admin.FullName, admin.PhoneNumber, hashPassword, admin.IsSuperAdmin)
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
		"phone_number": admin.PhoneNumber,
		"password":     admin.Password,
	})

}

func LoginAdmin(c *gin.Context) {

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
	var admin models.ShopOwnerLogin
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if !helpers.ValidatePhoneNumber(admin.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": errors.New("invalid phone number"),
		})
		return
	}

	// database - den telefon belgisi request - den gelen telefon belga den bolan maglumat cekilyar
	var id, password string
	if err := db.QueryRow(context.Background(), "SELECT id,password FROM admins WHERE phone_number = $1 AND deleted_at IS NULL", admin.PhoneNumber).Scan(&id, &password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status":  false,
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }
	// defer row.Close()

	// for row.Next() {
	// 	if err := row.Scan(&id, &password); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"status":  false,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}
	// }

	// eger request - den gelen telefon belgili admin database - de yok bolsa onda error response edilyar
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "this admin does not exist",
		})
		return
	}

	// eger shop_owner bar bolsa onda paroly dogry yazypdyrmy yazmandyrmy sol barlanyar
	credentialError := helpers.CheckPassword(admin.Password, password)
	if credentialError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid credentials",
		})
		return
	}

	// maglumatlar dogry bolsa auth ucin access_toke bilen resfresh_token generate edilyar
	accessTokenString, refreshTokenString, err := helpers.GenerateAccessTokenForAdmin(admin.PhoneNumber, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// front tarapa ugratmak ucin admin - in id - si boyunca maglumatlary get edilyar
	adm, err := GetAdminByID(id)
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

func GetAdminByID(id string) (models.ShopOwner, error) {
	db, err := config.ConnDB()
	if err != nil {
		return models.ShopOwner{}, err
	}
	defer db.Close()

	// parametrler edilip berilen id - boyunca database - den admin - in maglumatlary cekilyar
	var shopOwner models.ShopOwner
	rowShopOwner, err := db.Query(context.Background(), "SELECT full_name,phone_number FROM admins WHERE deleted_at IS NULL AND id = $1", id)
	if err != nil {
		return models.ShopOwner{}, err
	}
	defer rowShopOwner.Close()

	for rowShopOwner.Next() {
		if err := rowShopOwner.Scan(&shopOwner.FullName, &shopOwner.PhoneNumber); err != nil {
			return models.ShopOwner{}, err
		}
	}

	// eger parametrler edilip berilen id boyunca database - de maglumat yok bolsa error return edilyar
	if shopOwner.PhoneNumber == "" {
		return models.ShopOwner{}, errors.New("admin not found")
	}

	// hemme zat dogry bolsa admin - in maglumatlary return edilyar
	return shopOwner, nil

}

func UpdateAdmin(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request body - den admin - in maglumatlary alynyar
	var admin models.AdminUpdate
	if err := c.BindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// database - de request body - den gelen id bilen gabat gelyan admin barmy ya-da yokmy sol barlanyar
	// eger yok bolsa onda error return edilyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM admins WHERE id = $1 AND deleted_at IS NULL", admin.ID).Scan(&id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "admin not found",
		})
		return
	}

	if !helpers.ValidatePhoneNumber(admin.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": errors.New("invalid phone number"),
		})
		return
	}

	// eger admin database - de bar bolsa onda onun maglumatlary request body - dan gelen maglumatlar bilen update edilyar
	_, err = db.Exec(context.Background(), "UPDATE admins SET full_name = $1 , phone_number = $2 , is_super_admin = $3 WHERE id = $4", admin.FullName, admin.PhoneNumber, admin.IsSuperAdmin, id)
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
