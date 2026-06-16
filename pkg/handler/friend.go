package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type friendRequestInput struct {
	Username string `json:"username" binding:"required"`
}

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

func (h *Handler) searchUsers(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	prefix := c.Query("username")
	users, err := h.services.Authorization.SearchUsers(prefix)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"users": users})
}

func (h *Handler) getFriendList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	friends, err := h.services.Friendship.FriendList(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"friends": friends})
}

func (h *Handler) getFriendRequestList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	friends, err := h.services.Friendship.FriendRequestList(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"friends": friends})
}

func (h *Handler) getIncomingFriendRequests(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	friends, err := h.services.Friendship.IncomingFriendRequests(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"friends": friends})
}

func (h *Handler) getOutgoingFriendRequests(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	friends, err := h.services.Friendship.OutgoingFriendRequests(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"friends": friends})
}

func (h *Handler) sendFriendRequest(c *gin.Context) {
	var input friendRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	if err := h.services.Friendship.SendFriendRequest(userID, input.Username); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "friend request sent"})
}

func (h *Handler) acceptFriendRequest(c *gin.Context) {
	var input friendRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	if err := h.services.Friendship.AcceptFriendRequest(userID, input.Username); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "friend request accepted"})
}

func (h *Handler) rejectFriendRequest(c *gin.Context) {
	var input friendRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	if err := h.services.Friendship.RejectFriendRequest(userID, input.Username); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "friend request rejected"})
}

func (h *Handler) deleteFriend(c *gin.Context) {
	var input friendRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input")
		return
	}
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	if err := h.services.Friendship.DeleteFriend(userID, input.Username); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "friend deleted", "username": input.Username})
}
