package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateShop(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "salam",
	})

}
