package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ping(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": 0,
	})
}
