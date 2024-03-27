package controllers

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func RegisterShopOwner(c *gin.Context) {
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request - den gelen maglumatlar alynyar
	var shopOwner models.ShopOwner
	if err := c.BindJSON(&shopOwner); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// gelen maglumatlar barlanylyar
	if err := models.ValidateShopOwner(shopOwner, true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// parol hashlenyan
	hashPassword, err := helpers.HashPassword(shopOwner.Password)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// hemme zat yerbe yer bolsa maglumatlar shop_owners tablisa gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO shop_owners (full_name,phone_number,password,slug) VALUES ($1,$2,$3,$4)", shopOwner.FullName, shopOwner.PhoneNumber, hashPassword, slug.MakeLang(shopOwner.FullName, "en"))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
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
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request - den maglumatlar alynyar
	var shopOwner models.ShopOwnerLogin
	if err := c.BindJSON(&shopOwner); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if !helpers.ValidatePhoneNumber(shopOwner.PhoneNumber) {
		helpers.HandleError(c, 400, "invalid phone number")
		return
	}

	// database - den telefon belgisi request - den gelen telefon belga den bolan maglumat cekilyar
	var id, password string
	db.QueryRow(context.Background(), "SELECT id,password FROM shop_owners WHERE phone_number = $1 AND deleted_at IS NULL", shopOwner.PhoneNumber).Scan(&id, &password)

	// eger request - den gelen telefon belgili shop_owner database - de yok bolsa onda error response edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger shop_owner bar bolsa onda paroly dogry yazypdyrmy yazmandyrmy sol barlanyar
	credentialError := helpers.CheckPassword(shopOwner.Password, password)
	if credentialError != nil {
		helpers.HandleError(c, 400, "invalid credentials")
		return
	}

	// maglumatlar dogry bolsa auth ucin access_toke bilen resfresh_token generate edilyar
	accessTokenString, refreshTokenString, err := helpers.GenerateAccessTokenForAdmin(shopOwner.PhoneNumber, id, false)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// front tarapa ugratmak ucin shop_owner - in id - si boyunca maglumatlary get edilyar
	adm, err := GetShopOwnerByID(id)
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

func GetShopOwnerByID(id string) (models.ShopOwner, error) {
	db, err := config.ConnDB()
	if err != nil {
		return models.ShopOwner{}, err
	}
	defer db.Close()

	// parametrler edilip berilen id - boyunca database - den shop_owner - in maglumatlary cekilyar
	var shopOwner models.ShopOwner
	db.QueryRow(context.Background(), "SELECT full_name,phone_number FROM shop_owners WHERE deleted_at IS NULL AND id = $1", id).Scan(&shopOwner.FullName, &shopOwner.PhoneNumber)

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
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - den shop_owner - in maglumatlary alynyar
	var shopOwner models.ShopOwner
	if err := c.BindJSON(&shopOwner); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateShopOwner(shopOwner, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger shop_owner database - de bar bolsa onda onun maglumatlary request body - dan gelen maglumatlar bilen update edilyar
	_, err = db.Exec(context.Background(), "UPDATE shop_owners SET full_name = $1 , phone_number = $2 , slug = $3 WHERE id = $4", shopOwner.FullName, shopOwner.PhoneNumber, slug.MakeLang(shopOwner.FullName, "en"), shopOwner.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func GetShopOwners(c *gin.Context) {
	var requestQuery helpers.StandartQuery
	var count uint
	var shopOwners []models.ShopOwner

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

	// request query - den status - a gora shop owner - leryn sanyny almak ucin query yazylyar
	queryCount := `SELECT COUNT(id) FROM shop_owners WHERE deleted_at IS NULL`
	if requestQuery.IsDeleted {
		queryCount = `SELECT COUNT(id) FROM shop_owners WHERE deleted_at IS NOT NULL`
	}
	if err := db.QueryRow(context.Background(), queryCount).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora shop owner - lary almak ucin query yazylyar
	rowQuery := `SELECT full_name,phone_number FROM shop_owners WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if requestQuery.IsDeleted {
		rowQuery = `SELECT full_name,phone_number FROM shop_owners WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}
	// databae - den request - den gelen limit we page boyunca limitlap shop_owner - ler alynyar
	rowsShopOwner, err := db.Query(context.Background(), rowQuery, requestQuery.Limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsShopOwner.Close()

	for rowsShopOwner.Next() {
		var shopOwner models.ShopOwner
		if err := rowsShopOwner.Scan(&shopOwner.FullName, &shopOwner.PhoneNumber); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shopOwners = append(shopOwners, shopOwner)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      true,
		"shop_owners": shopOwners,
		"total":       count,
	})
}

func GetShopOwner(c *gin.Context) {
	adm, err := GetShopOwnerByID(c.Param("id"))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"shop_owner": adm,
	})
}

func DeleteShopOwnerByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den shop_owner id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("shop_owners", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa shop_owner - in we sol shop_owner - a degisli shops - laryn deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_shop_owner($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreShopOwnerByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den shop id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("shop_owners", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa maglumat restore edilyar
	_, err = db.Exec(context.Background(), "CALL restore_shop_owner($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyShopOwnerByID(c *gin.Context) {
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
	rowShopOwner, err := db.Query(context.Background(), "SELECT so.id,s.image FROM shop_owners so INNER JOIN shops s ON s.shop_owner_id = so.id WHERE so.id = $1 AND so.deleted_at IS NOT NULL", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	var id string
	var shopImages []string
	for rowShopOwner.Next() {
		var shopImage string
		if err := rowShopOwner.Scan(&id, &shopImage); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		shopImages = append(shopImages, shopImage)
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// eger maglumat bar bolsa sonda shop - yn suratlaryny papkadan pozulyar
	for _, v := range shopImages {
		if err := os.Remove(helpers.ServerPath + v); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if err := os.Remove(helpers.ServerPath + "assets/" + v); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	// shop - yn suraty pozulandan sonra database - den shop - lar pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shops WHERE shop_owner_id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// shop - lar pozulandan sonra shop_owner database - den pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM shop_owners WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func UpdateShopOwnerPassword(c *gin.Context) {
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

	if err := helpers.ValidateRecordByID("shop_owners", admin.ID, "NULL", db); err != nil {
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
	_, err = db.Exec(context.Background(), "UPDATE shop_owners SET password = $1 WHERE id = $2", hashPassword, admin.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "password of shop owner successfuly updated",
	})
}
