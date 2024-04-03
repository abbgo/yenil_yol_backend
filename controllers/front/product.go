package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProductByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametrden shop id alynyar
	productID := c.Param("id")

	// request - dan gelen id boyunca haryt alynyar
	var product models.Product
	db.QueryRow(context.Background(), "SELECT id,name_tm,name_ru,price,old_price,brend_id FROM products WHERE id=$1 AND deleted_at IS NULL", productID).Scan(&product.ID, &product.NameTM, &product.NameRU, &product.Price, &product.OldPrice, &product.BrendID)
	if product.ID == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// harydyn brendi alynyar
	db.QueryRow(context.Background(), "SELECT name,image FROM brends WHERE id=$1 AND deleted_at IS NULL", product.BrendID).Scan(&product.Brend.Name, &product.Brend.Image)

	// sonra harydyn renkleri we ona degisli suratlar alynyar
	var productColor models.ProductColor
	rowsColor, err := db.Query(context.Background(), "SELECT id FROM product_colors WHERE product_id=$1", product.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsColor.Close()
	for rowsColor.Next() {
		if err := rowsColor.Scan(&productColor.ID); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		// sonra renke degisli suratlar alynyar
		rowsImage, err := db.Query(context.Background(), "SELECT image FROM product_images WHERE product_color_id=$1 AND deleted_at IS NULL", productColor.ID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		defer rowsColor.Close()
		for rowsImage.Next() {
			var image string
			if err := rowsImage.Scan(&image); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}
			productColor.Images = append(productColor.Images, image)
		}

		product.ProductColors = append(product.ProductColors, productColor)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"product": product,
	})
}
