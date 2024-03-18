package helpers

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FileUpload(fileName, path, fileType string, context *gin.Context) (string, error) {

	file, err := context.FormFile(fileName)
	if err != nil {
		return "", err
	}

	extensionFile := filepath.Ext(file.Filename)

	var newFileName string

	// VALIDATE IMAGE
	if fileType == "image" {
		if extensionFile != ".jpg" && extensionFile != ".JPG" && extensionFile != ".jpeg" && extensionFile != ".JPEG" && extensionFile != ".png" && extensionFile != ".PNG" && extensionFile != ".gif" && extensionFile != ".GIF" && extensionFile != ".svg" && extensionFile != ".SVG" && extensionFile != ".WEBP" && extensionFile != ".webp" {
			return "", errors.New("the file must be an image")
		}
		newFileName = uuid.New().String() + extensionFile

	}

	_, err = os.Stat(ServerPath + "uploads/" + path)
	if err != nil {
		if err := os.MkdirAll(ServerPath+"uploads/"+path, os.ModePerm); err != nil {
			return "", err
		}
	}
	if err := context.SaveUploadedFile(file, ServerPath+"uploads/"+path+"/"+newFileName); err != nil {
		return "", err
	}

	_, err = os.Stat(ServerPath + "assets/uploads/" + path)
	if err != nil {
		if err := os.MkdirAll(ServerPath+"assets/uploads/"+path, os.ModePerm); err != nil {
			return "", err
		}
	}

	src, err := imaging.Open(ServerPath + "uploads/" + path + "/" + newFileName)
	if err != nil {
		return "", err
	}

	src = imaging.Resize(src, 200, 0, imaging.Lanczos)

	err = imaging.Save(src, ServerPath+"assets/uploads/"+path+"/"+newFileName)
	if err != nil {
		return "", err
	}

	return "uploads/" + path + "/" + newFileName, nil

}
