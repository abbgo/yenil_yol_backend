package controllers

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterAdmin(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request - den gelen maglumatlar alynyar
	var admin models.Admin
	if err := c.BindJSON(&admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateAdmin(admin.PhoneNumber, "", true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// parol hashlenyan
	hashPassword, err := helpers.HashPassword(admin.Password)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// hemme zat yerbe yer bolsa maglumatlar admins tablisa gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO admins (full_name,phone_number,password,is_super_admin) VALUES ($1,$2,$3,$4)", admin.FullName, admin.PhoneNumber, hashPassword, admin.IsSuperAdmin)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"phone_number": admin.PhoneNumber,
		"full_name":    admin.FullName,
	})

}

func LoginAdmin(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request - den maglumatlar alynyar
	var admin models.ShopOwnerLogin
	if err := c.BindJSON(&admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if !helpers.ValidatePhoneNumber(admin.PhoneNumber) {
		helpers.HandleError(c, 400, "invalid phone number")
		return
	}

	// database - den telefon belgisi request - den gelen telefon belga den bolan maglumat cekilyar
	var id, password string
	var is_super_admin bool
	if err := db.QueryRow(context.Background(), "SELECT id,password,is_super_admin FROM admins WHERE phone_number = $1 AND deleted_at IS NULL", admin.PhoneNumber).Scan(&id, &password, &is_super_admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger request - den gelen telefon belgili admin database - de yok bolsa onda error response edilyar
	if id == "" {
		helpers.HandleError(c, 404, "this admin does not exist")
		return
	}

	// eger admin bar bolsa onda paroly dogry yazypdyrmy yazmandyrmy sol barlanyar
	credentialError := helpers.CheckPassword(admin.Password, password)
	if credentialError != nil {
		helpers.HandleError(c, 400, "invalid credentials")
		return
	}

	// maglumatlar dogry bolsa auth ucin access_toke bilen resfresh_token generate edilyar
	accessTokenString, refreshTokenString, err := helpers.GenerateAccessTokenForAdmin(admin.PhoneNumber, id, is_super_admin)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// front tarapa ugratmak ucin admin - in id - si boyunca maglumatlary get edilyar
	adm, err := GetAdminByID(id)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
		"admin":         adm,
	})

}

func GetAdminByID(id string) (models.Admin, error) {
	db, err := config.ConnDB()
	if err != nil {
		return models.Admin{}, err
	}
	defer db.Close()

	// parametrler edilip berilen id - boyunca database - den admin - in maglumatlary cekilyar
	var admin models.Admin
	rowAdmin, err := db.Query(context.Background(), "SELECT full_name,phone_number,is_super_admin FROM admins WHERE deleted_at IS NULL AND id = $1", id)
	if err != nil {
		return models.Admin{}, err
	}
	defer rowAdmin.Close()

	for rowAdmin.Next() {
		if err := rowAdmin.Scan(&admin.FullName, &admin.PhoneNumber, &admin.IsSuperAdmin); err != nil {
			return models.Admin{}, err
		}
	}

	// eger parametrler edilip berilen id boyunca database - de maglumat yok bolsa error return edilyar
	if admin.PhoneNumber == "" {
		return models.Admin{}, errors.New("admin not found")
	}

	// hemme zat dogry bolsa admin - in maglumatlary return edilyar
	return admin, nil

}

func UpdateAdmin(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - den admin - in maglumatlary alynyar
	var admin models.AdminUpdate
	if err := c.BindJSON(&admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateAdmin(admin.PhoneNumber, admin.ID, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger admin database - de bar bolsa onda onun maglumatlary request body - dan gelen maglumatlar bilen update edilyar
	_, err = db.Exec(context.Background(), "UPDATE admins SET full_name = $1 , phone_number = $2 , is_super_admin = $3 WHERE id = $4", admin.FullName, admin.PhoneNumber, admin.IsSuperAdmin, admin.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}

func GetAdmin(c *gin.Context) {

	adminID, hasID := c.Get("admin_id")
	if !hasID {
		helpers.HandleError(c, 400, "adminID is required")
		return
	}

	var ok bool
	admin_id, ok := adminID.(string)
	if !ok {
		helpers.HandleError(c, 400, "adminID must be string")
		return
	}

	adm, err := GetAdminByID(admin_id)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"admin": adm,
	})

}

func GetAdmins(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den limit alynyar
	limitStr := c.Query("limit")
	if limitStr == "" {
		helpers.HandleError(c, 400, "limit is required")
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// // request parametr - den page alynyar
	pageStr := c.Query("page")
	if pageStr == "" {
		helpers.HandleError(c, 400, "page is required")
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// database - den admin - lerin sany alynyar
	countOfAdmins := 0
	if err := db.QueryRow(context.Background(), "SELECT COUNT(id) FROM admins WHERE deleted_at IS NULL").Scan(&countOfAdmins); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// databae - den request - den gelen limit we page boyunca limitlap admin - ler alynyar
	var admins []models.Admin
	rowsAdmin, err := db.Query(context.Background(), "SELECT full_name,phone_number,is_super_admin FROM admins WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsAdmin.Close()

	for rowsAdmin.Next() {
		var admin models.Admin
		if err := rowsAdmin.Scan(&admin.FullName, &admin.PhoneNumber, &admin.IsSuperAdmin); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		admins = append(admins, admin)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"admins": admins,
		"total":  countOfAdmins,
	})

}

func DeleteAdminByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den admin id alynyar
	ID := c.Param("id")

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM admins WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa admin - in  deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE admins SET deleted_at=NOW() WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}

func RestoreAdminByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den admin id alynyar
	ID := c.Param("id")

	// alynan id den bolan maglumat database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM admins WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database sol id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa maglumat restore edilyar
	_, err = db.Exec(context.Background(), "UPDATE admins SET deleted_at=NULL WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})

}

func DeletePermanentlyAdminByID(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den admin id alynyar
	ID := c.Param("id")

	// database - de gelen id degisli maglumat barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM admins WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// pozulandan sonra admin database - den pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM admins WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})

}

func UpdateAdminPassword(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - den maglumatlar alynyar
	var admin models.AdminUpdatePassword
	if err := c.BindJSON(&admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var admin_id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM admins WHERE id = $1 AND deleted_at IS NULL", admin.ID).Scan(&admin_id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger gelen id den bolan maglumat database - de yok bolsa error return edilyar
	if admin_id == "" {
		helpers.HandleError(c, 404, "admin not found")
		return
	}

	// maglumat bar bolsa admin - in taze paroly hashlenyar
	hashPassword, err := helpers.HashPassword(admin.Password)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// taze parol kone parol bilen calsylyar
	_, err = db.Exec(context.Background(), "UPDATE admins SET password = $1 WHERE id = $2", hashPassword, admin.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "password of admin successfuly updated",
	})

}
