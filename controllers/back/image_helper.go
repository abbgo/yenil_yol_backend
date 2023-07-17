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
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	var path, file_name string
	imageType := c.Query("image_type")

	oldImage := c.PostForm("old_path")
	if oldImage != "" {
		if err := os.Remove(helpers.ServerPath + oldImage); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}

		_, err := db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", oldImage)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
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
	case "product":
		path = "product"
		file_name = "image"
	case "setting":
		path = "setting"
		file_name = "image"
	case "category":
		path = "category"
		file_name = "image"
	case "brend":
		path = "brend"
		file_name = "image"
	case "shop":
		path = "shop"
		file_name = "image"
	case "page":
		path = "page"
		file_name = "image"
	default:
		helpers.HandleError(c, 400, "invalid image")
		return
	}

	image, err := helpers.FileUpload(file_name, path, "image", c)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	_, err = db.Exec(context.Background(), "INSERT INTO helper_images (image) VALUES ($1)", image)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"image":  image,
	})

}

type DeleteImg struct {
	Image string `json:"image"`
}

func DeleteImage(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	var image DeleteImg
	if err := c.Bind(&image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	if image.Image == "" {
		helpers.HandleError(c, 400, "path of image is required")
		return
	}

	var helperImageID string
	if err := db.QueryRow(context.Background(), "SELECT id FROM helper_images WHERE image = $1 AND deleted_at IS NULL", image.Image).Scan(&helperImageID); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	if helperImageID != "" {
		_, err := db.Exec(context.Background(), "DELETE FROM helper_images WHERE id = $1", helperImageID)
		if err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
	}

	if err := os.Remove(helpers.ServerPath + image.Image); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "image successfully deleted",
	})

}
