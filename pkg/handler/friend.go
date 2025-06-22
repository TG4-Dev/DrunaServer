package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getFriendList(c *gin.Context) {
	id, _ := c.Get(userCtx)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) sendFriendRequest(c *gin.Context) {

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
