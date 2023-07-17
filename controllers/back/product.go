package controllers

import (
	"context"
	"fmt"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func CreateProduct(c *gin.Context) {

	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - dan gelen maglumatlar alynyar
	var product models.Product
	if err := c.BindJSON(&product); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger maglumatlar dogry bolsa onda products tablisa maglumatlar gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO products (name_tm,name_ru,image,price,old_price,status,color_name_tm,color_name_ru,gender_name_tm,gender_name_ru,code,slug_tm,slug_ru,shop_id,category_id,brend_id) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)", product.NameTM, product.NameRU, product.Image, product.Price, product.OldPrice, product.Status, product.ColorNameTM, product.ColorNameRU, product.GenderNameTM, product.GenderNameRU, product.Code, slug.MakeLang(product.NameTM, "en"), slug.MakeLang(product.NameRU, "en"), product.ShopID, product.CategoryID, product.BrendID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		fmt.Println("yalnys 1")
		return
	}

	// brend - yn maglumatlary gosulandan sonra helper_images tablisa brend ucin gosulan surat pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", product.Image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully added",
	})

}
