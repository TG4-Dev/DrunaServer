package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EventAddHandler(c *gin.Context) { //GET
	c.JSON(http.StatusOK, gin.H{
		"message": "adding event",
	})
}

func EventGetFreeTime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mesage": "calculating time",
	})
}
