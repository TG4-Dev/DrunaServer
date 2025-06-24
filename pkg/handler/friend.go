package handler

import (
	"druna_server/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getFriendList(c *gin.Context) {
	type FriendList struct {
		SourceID int `json:"SourceID" binding:"required"`
	}

	var input FriendList
	var friends []model.FriendInfo

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	friends, err := h.services.Friendship.FriendList(input.SourceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "getting friend list",
		"friends": friends,
	})
}

func (h *Handler) sendFriendRequest(c *gin.Context) {
	type FriendRequest struct {
		SourceID int    `json:"SourceID" binding:"required"`
		Username string `json:"username" binding:"required"`
	}
	var input FriendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.services.Friendship.FriendRequest(input.SourceID, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request sent",
	})
}

func (h *Handler) acceptFriendRequest(c *gin.Context) {
}
func (h *Handler) deleteFriend(c *gin.Context) {
	id, _ := c.Get(userCtx)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getFriendRequestList(c *gin.Context) {
	id, _ := c.Get(userCtx)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
