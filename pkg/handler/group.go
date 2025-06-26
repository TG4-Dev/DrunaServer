package handler

import (
	"druna_server/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createGroup(c *gin.Context) {

	var input model.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.services.Group.CreateGroup(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "group created",
		"groupId": id,
	})
}
