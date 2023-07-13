package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func CreateBrend(c *gin.Context) {

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
	var brend models.Brend
	if err := c.BindJSON(&brend); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// eger request body - dan gelen surat bos bolsa database surata derek nil gosmaly
	var image interface{}
	if brend.Image == "" {
		image = nil
	} else {
		image = brend.Image
	}

	// eger maglumatlar dogry bolsa onda brends tablisa maglumatlar gosulyar we gosulandan son gosulan maglumatyn id - si return edilyar
	_, err = db.Exec(context.Background(), "INSERT INTO brends (name,image,slug) VALUES ($1,$2,$3)", brend.Name, image, slug.MakeLang(brend.Name, "en"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa shop ucin gosulan surat pozulyar
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
