package controllers

import (
	"context"
	"database/sql"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CreateSetting(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var setting models.Setting
	if err := c.BindJSON(&setting); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateSetting(setting.PhoneNumber, setting.Email); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda brends tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	_, err = db.Exec(context.Background(), "INSERT INTO settings (logo,favicon,email,phone_number) VALUES ($1,$2,$3,$4)", setting.Logo, setting.Favicon, setting.Email, setting.PhoneNumber)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// setting - yn maglumatlary gosulandan sonra helper_images tablisa setting ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1 OR image =$2", setting.Logo, setting.Favicon)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})

}

func UpdateSetting(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var setting models.SettingUpdate
	if err := c.BindJSON(&setting); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request body - da gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var oldLogo, oldFavicon sql.NullString
	if err := db.QueryRow(context.Background(), "SELECT logo,favicon FROM settings WHERE deleted_at IS NULL").Scan(&oldLogo, &oldFavicon); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de sol maglumat yok bolsa onda error return edilyar
	if oldFavicon.String == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// bu yerde logo ucin newLogo atly uytgeyan ululyk doredilyar
	// eger request body - dan logo gelmese onda onki logo uytgedilmeyar diymek bolyar
	// sonun ucin onki logany goyyas , eger request body - dan logo gelen bolsa
	// onda taze logany kone logo bilen calysyas
	var newLogo string
	if setting.Logo == "" {
		newLogo = oldLogo.String
	} else {
		// sonra helper_images tablisadan logo ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", setting.Logo)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if oldLogo.String != "" {
			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + oldLogo.String); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		newLogo = setting.Logo
	}

	// bu yerde favicon ucin newFavicon atly uytgeyan ululyk doredilyar
	// eger request body - dan favicon gelmese onda onki favicon uytgedilmeyar diymek bolyar
	// sonun ucin onki favicony goyyas , eger request body - dan favicon gelen bolsa
	// onda taze favicony kone favicon bilen calysyas
	var newFavicon string
	if setting.Favicon == "" {
		newFavicon = oldFavicon.String
	} else {
		// sonra helper_images tablisadan favicon ucin gosulan surat pozulyar
		_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", setting.Favicon)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		if oldFavicon.String != "" {
			// surat papkadan pozulyar
			if err := os.Remove(helpers.ServerPath + oldFavicon.String); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
		}
		newFavicon = setting.Favicon
	}

	// database - daki maglumatlary request body - dan gelen maglumatlar bilen calysyas
	_, err = db.Exec(context.Background(), "UPDATE settings SET logo=$1 , favicon=$2 , email=$3 , phone_number=$4", newLogo, newFavicon, setting.Email, setting.PhoneNumber)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})

}
