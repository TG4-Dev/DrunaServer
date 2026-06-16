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

// @Summary List friends
// @Security ApiKeyAuth
// @Tags friends
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/friends/list [get]
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

	c.JSON(http.StatusOK, gin.H{
		"friends": friends,
	})
}

// @Summary List all pending friend requests
// @Security ApiKeyAuth
// @Tags friends
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/request-list [get]
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

	c.JSON(http.StatusOK, gin.H{
		"friends": friends,
	})
}

// @Summary List incoming friend requests
// @Security ApiKeyAuth
// @Tags friends
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/requests/incoming [get]
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

	c.JSON(http.StatusOK, gin.H{
		"friends": friends,
	})
}

// @Summary List outgoing friend requests
// @Security ApiKeyAuth
// @Tags friends
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/requests/outgoing [get]
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

	c.JSON(http.StatusOK, gin.H{
		"friends": friends,
	})
}

// @Summary Send friend request
// @Security ApiKeyAuth
// @Tags friends
// @Accept json
// @Produce json
// @Param input body model.FriendRequestDoc true "username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/request [post]
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

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request sent",
	})
}

// @Summary Accept friend request
// @Security ApiKeyAuth
// @Tags friends
// @Accept json
// @Produce json
// @Param input body model.FriendRequestDoc true "username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/accept [post]
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

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request accepted",
	})
}

// @Summary Reject friend request
// @Security ApiKeyAuth
// @Tags friends
// @Accept json
// @Produce json
// @Param input body model.FriendRequestDoc true "username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/reject [post]
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

	c.JSON(http.StatusOK, gin.H{
		"message": "friend request rejected",
	})
}

// @Summary Delete friend
// @Security ApiKeyAuth
// @Tags friends
// @Accept json
// @Produce json
// @Param input body model.FriendRequestDoc true "username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/friends/ [delete]
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

	c.JSON(http.StatusOK, gin.H{
		"message": "friend deleted",
		"username": input.Username,
	})
}
