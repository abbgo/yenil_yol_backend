package helpers

import (
	"errors"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FileUpload(fileName, path, fileType string, context *gin.Context /* resizedSize int */) (string, error) {

	file, err := context.FormFile(fileName)
	if err != nil {
		return "", err
	}

	extensionFile := filepath.Ext(file.Filename)

	var newFileName string

	// VALIDATE IMAGE
	// if fileType == "image" {
	// 	if extensionFile != ".jpg" && extensionFile != ".JPG" && extensionFile != ".jpeg" && extensionFile != ".JPEG" && extensionFile != ".png" && extensionFile != ".PNG" && extensionFile != ".gif" && extensionFile != ".GIF" && extensionFile != ".svg" && extensionFile != ".SVG" && extensionFile != ".WEBP" && extensionFile != ".webp" {
	// 		return "", errors.New("the file must be an image")
	// 	}
	// 	newFileName = uuid.New().String() + extensionFile
	// }
	if fileType == "image" {
		if extensionFile != ".jpg" && extensionFile != ".JPG" && extensionFile != ".jpeg" && extensionFile != ".JPEG" {
			return "", errors.New("the image must be .jpg format")
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

	// Resize Image
	// if resizedSize != 0 {
	// 	if err := ResizeImage(path, newFileName, resizedSize); err != nil {
	// 		return "", err
	// 	}
	// }

	if err := CompressImage(path, newFileName); err != nil {
		return "", err
	}

	return "uploads/" + path + "/" + newFileName, nil

}

func CompressImage(path string, newFileName string) error {
	_, err := os.Stat(ServerPath + "assets/uploads/" + path)
	if err != nil {
		if err := os.MkdirAll(ServerPath+"assets/uploads/"+path, os.ModePerm); err != nil {
			return err
		}
	}

	f, err := os.Open(ServerPath + "uploads/" + path + "/" + newFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	out, err := os.Create(ServerPath + "assets/uploads/" + path + "/" + newFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	opts := &jpeg.Options{Quality: 80} // Adjust quality (1-100)
	err = jpeg.Encode(out, img, opts)
	return err
}

// func ResizeImage(path string, newFileName string, imageWidth int) error {
// 	_, err := os.Stat(ServerPath + "assets/uploads/" + path)
// 	if err != nil {
// 		if err := os.MkdirAll(ServerPath+"assets/uploads/"+path, os.ModePerm); err != nil {
// 			return err
// 		}
// 	}

// 	src, err := imaging.Open(ServerPath + "uploads/" + path + "/" + newFileName)
// 	if err != nil {
// 		return err
// 	}

// 	src = imaging.Resize(src, imageWidth, 0, imaging.Lanczos)

// 	err = imaging.Save(src, ServerPath+"assets/uploads/"+path+"/"+newFileName)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
