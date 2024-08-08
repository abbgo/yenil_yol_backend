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

	if err := models.ValidateAdmin(admin, true); err != nil {
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
	var admin models.Login
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
	db.QueryRow(context.Background(), "SELECT id,password,is_super_admin FROM admins WHERE phone_number = $1 AND deleted_at IS NULL", admin.PhoneNumber).Scan(&id, &password, &is_super_admin)

	// eger request - den gelen telefon belgili admin database - de yok bolsa onda error response edilyar
	if id == "" {
		helpers.HandleError(c, 404, "admin not found")
		return
	}

	// eger admin bar bolsa onda paroly dogry yazypdyrmy yazmandyrmy sol barlanyar
	credentialError := helpers.CheckPassword(admin.Password, password)
	if credentialError != nil {
		helpers.HandleError(c, 400, "invalid credentials")
		return
	}

	// maglumatlar dogry bolsa auth ucin access_toke bilen resfresh_token generate edilyar
	accessTokenString /* refreshTokenString,*/, err := helpers.GenerateAccessTokenForAdmin(admin.PhoneNumber, id, is_super_admin)
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
		"access_token": accessTokenString,
		// "refresh_token": refreshTokenString,
		"admin":  adm,
		"status": true,
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
	db.QueryRow(context.Background(), "SELECT full_name,phone_number,is_super_admin FROM admins WHERE deleted_at IS NULL AND id = $1", id).Scan(&admin.FullName, &admin.PhoneNumber, &admin.IsSuperAdmin)

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
	var admin models.Admin
	if err := c.BindJSON(&admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateAdmin(admin, false); err != nil {
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
	var requestQuery helpers.StandartQuery
	var countOfAdmins uint
	var admins []models.Admin

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

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

	// database - den admin - lerin sany alynyar
	queryCount := `SELECT COUNT(id) FROM admins WHERE deleted_at IS NULL`
	if requestQuery.IsDeleted {
		queryCount = `SELECT COUNT(id) FROM admins WHERE deleted_at IS NOT NULL`
	}
	if err := db.QueryRow(context.Background(), queryCount).Scan(&countOfAdmins); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// databae - den request - den gelen limit we page boyunca limitlap admin - ler alynyar
	rowQuery := `SELECT full_name,phone_number,is_super_admin FROM admins WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if requestQuery.IsDeleted {
		rowQuery = `SELECT full_name,phone_number,is_super_admin FROM admins WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}
	rowsAdmin, err := db.Query(context.Background(), rowQuery, requestQuery.Limit, offset)
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
	if err := helpers.ValidateRecordByID("admins", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
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
	if err := helpers.ValidateRecordByID("admins", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
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
	if err := helpers.ValidateRecordByID("admins", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
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
	var admin models.UpdatePassword
	if err := c.BindJSON(&admin); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("admins", admin.ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
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
