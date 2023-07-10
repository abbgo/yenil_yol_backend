package helpers

import (
	"errors"
	"os"
	"path/filepath"

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

	} /*else if fileType == "excel" {
		if extensionFile != ".xlsx" && extensionFile != ".xlsm" && extensionFile != ".xlsb" && extensionFile != ".xltx" {
			return "", errors.New("the file must be an excel")
		}
		newFileName = "products" + extensionFile
	} else {
		return "", errors.New("invalid file type")
	}*/

	_, err = os.Stat(ServerPath + "uploads/" + path)
	if err != nil {
		if err := os.MkdirAll(ServerPath+"uploads/"+path, os.ModePerm); err != nil {
			return "", err
		}
	}
	if err := context.SaveUploadedFile(file, ServerPath+"uploads/"+path+"/"+newFileName); err != nil {
		return "", err
	}

	return "uploads/" + path + "/" + newFileName, nil

}
