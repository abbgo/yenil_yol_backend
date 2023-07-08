package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterShopOwner(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "RegisterShopOwner controller",
	})

}
