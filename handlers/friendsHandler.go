package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FriendListHandler(c *gin.Context) { //GET
	c.JSON(http.StatusOK, gin.H{
		"message": "getting friend list",
	})
}

func FriendRequestHandler(c *gin.Context) { //POST
	c.JSON(http.StatusOK, gin.H{
		"message": "adding friend",
	})
}
