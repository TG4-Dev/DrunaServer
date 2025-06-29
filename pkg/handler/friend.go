package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUserIdFromToken(c *gin.Context) int {
	id, ok := c.Get(userCtx)

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0
	}

	userID, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0
	}

	return userID
}

func (h *Handler) getFriendRequestList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	friends, err := h.services.Friendship.FriendRequestList(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "getting friend list",
		"friends": friends,
	})
}

func (h *Handler) getFriendList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	friends, err := h.services.Friendship.FriendList(userID)
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
		Username string `json:"username" binding:"required"`
	}

	var input FriendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	err := h.services.Friendship.SendFriendRequest(userID, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request sent",
	})
}

func (h *Handler) acceptFriendRequest(c *gin.Context) {
	type FriendRequest struct {
		Username string `json:"username" binding:"required"`
	}

	var input FriendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	err := h.services.Friendship.AcceptFriendRequest(userID, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request accepted",
	})

}

func (h *Handler) rejectFriendRequest(c *gin.Context) {
	type FriendRequest struct {
		Username string `json:"username" binding:"required"`
	}

	var input FriendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	err := h.services.Friendship.RejectFriendRequest(userID, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request rejected",
	})
}

func (h *Handler) deleteFriend(c *gin.Context) {
	type FriendRequest struct {
		Username string `json:"username" binding:"required"`
	}

	var input FriendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	err := h.services.Friendship.DeleteFriend(userID, input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		input.Username: "Deleted",
	})
}
