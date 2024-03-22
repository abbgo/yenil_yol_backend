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
	var resizedSize int

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
	case "product":
		path = "product"
		file_name = "image"
		resizedSize = 200
	case "setting":
		path = "setting"
		file_name = "image"
		resizedSize = 0
	case "category":
		path = "category"
		file_name = "image"
		resizedSize = 200
	case "brend":
		path = "brend"
		file_name = "image"
		resizedSize = 200
	case "shop":
		path = "shop"
		file_name = "image"
		resizedSize = 230
	default:
		helpers.HandleError(c, 400, "invalid image")
		return
	}

	image, err := helpers.FileUpload(file_name, path, "image", c, resizedSize)
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
	Image string `json:"image" binding:"required"`
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

	var helperImageID string
	db.QueryRow(context.Background(), "SELECT id FROM helper_images WHERE image = $1 AND deleted_at IS NULL", image.Image).Scan(&helperImageID)
	if helperImageID == "" {
		helpers.HandleError(c, 404, "record not found")
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

func DeleteImageFromDB(image string) error {
	db, err := config.ConnDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(context.Background(), "DELETE FROM helper_images WHERE image = $1", image)
	if err != nil {
		return err
	}

	return nil
}
