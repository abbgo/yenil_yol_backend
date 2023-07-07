package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHeaderData(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "salam",
	})

}
