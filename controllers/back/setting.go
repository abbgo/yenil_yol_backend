package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSetting(c *gin.Context) {

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
	var setting models.Setting
	if err := c.BindJSON(&setting); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := models.ValidateSetting(setting.PhoneNumber, setting.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// eger maglumatlar dogry bolsa onda brends tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	_, err = db.Exec(context.Background(), "INSERT INTO settings (logo,favicon,email,phone_number) VALUES ($1,$2,$3,$4)", setting.Logo, setting.Favicon, setting.Email, setting.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// setting - yn maglumatlary gosulandan sonra helper_images tablisa setting ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1 OR image =$2", setting.Logo, setting.Favicon)
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
