package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create group
// @Security ApiKeyAuth
// @Tags groups
// @Accept json
// @Produce json
// @Param input body model.Group true "group info"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/groups/create [post]
func (h *Handler) createGroup(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	var input model.Group
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input.OwnerID = userID

	id, err := h.services.Group.CreateGroup(input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "group created",
		"groupId": id,
	})
}

// @Summary List user groups
// @Security ApiKeyAuth
// @Tags groups
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/groups/list [get]
func (h *Handler) listGroups(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	groups, err := h.services.Group.ListGroups(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"groups": groups,
	})
}

// @Summary Get group details
// @Security ApiKeyAuth
// @Tags groups
// @Produce json
// @Param id path int true "Group ID"
// @Success 200 {object} model.GroupDetails
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Router /api/groups/{id} [get]
func (h *Handler) getGroup(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}

	details, err := h.services.Group.GetGroupDetails(groupID, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, details)
}

// @Summary Add group member
// @Security ApiKeyAuth
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param input body model.AddGroupMemberDoc true "member username"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/groups/{id}/members [post]
func (h *Handler) addGroupMember(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}

	var input model.AddGroupMemberDoc
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Group.AddGroupMember(groupID, userID, input.Username); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "member added",
	})
}
