package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterShopOwner(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	var shopOwner models.ShopOwner
	if err := c.BindJSON(&shopOwner); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := models.ValidateRegisterShopOwner(shopOwner.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	hashPassword, err := helpers.HashPassword(shopOwner.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}

	_, err = db.Exec(context.Background(), "INSERT INTO shop_owners (name,phone_number,password) VALUES ($1,$2,$3)", shopOwner.Name, shopOwner.PhoneNumber, hashPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"error":   "1",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"phone_number": shopOwner.PhoneNumber,
		"password":     shopOwner.Password,
		// "admin_type":   admin.Type,
	})

}
