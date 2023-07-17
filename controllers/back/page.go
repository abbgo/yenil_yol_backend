package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePage(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var page models.Page
	if err := c.BindJSON(&page); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if page.Image == "" {
		image = nil
	} else {
		image = page.Image
	}

	// eger maglumatlar dogry bolsa onda pages tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO pages (name,image) VALUES ($1,$2)", page.Name, image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa page ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})

}
