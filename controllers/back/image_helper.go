package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddOrUpdateImage(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	var path, file_name string
	imageType := c.Query("image_type")

	oldImage := c.PostForm("old_path")
	if oldImage != "" {
		if err := os.Remove(helpers.ServerPath + oldImage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}

		_, err := db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", oldImage)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": err.Error(),
			})
			return
		}
	}

	switch imageType {
	// case "product":
	// 	fileName := c.Query("type")
	// 	if fileName != "main_image" && fileName != "image" {
	// 		c.JSON(http.StatusNotFound, gin.H{
	// 			"status":  false,
	// 			"message": "invalid file name",
	// 		})
	// 		return
	// 	}
	// 	path = "product/" + fileName
	// 	file_name = fileName
	// case "category":
	// 	path = "category"
	// 	file_name = "image"
	// case "brend":
	// 	path = "brend"
	// 	file_name = "image"
	// case "language":
	// 	path = "language"
	// 	file_name = "image"
	// case "banner":
	// 	path = "banner"
	// 	file_name = "image"
	case "shop":
		path = "shop"
		file_name = "image"
	default:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "invalid image",
		})
		return
	}

	image, err := helpers.FileUpload(file_name, path, "image", c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	_, err = db.Exec(context.Background(), "INSERT INTO helper_images (image) VALUES ($1)", image)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	// defer func() {
	// 	if err := result.Close(); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"status":  false,
	// 			"message": err.Error(),
	// 		})
	// 		return
	// 	}
	// }()

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"image":  image,
	})

}
