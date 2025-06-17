package handlers

import (
	"net/http"

	"BlobbyServer/pkg/services"

	"github.com/gin-gonic/gin"
)

func FriendListHandler(c *gin.Context) { //GET
	c.JSON(http.StatusOK, gin.H{
		"message": "getting friend list",
	})
}

func FriendRequestHandler(c *gin.Context) { //POST

	type FriendRequest struct {
		SourceID int    `json:"SourceID" binding:"required"`
		Username string `json:"username" binding:"required"`
	}

	var input FriendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := services.FriendService.FriendRequest(input.SourceID, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "adding friend",
	})
}
