package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
	Success(c, http.StatusOK, gin.H{"message": "group created", "groupId": id})
}

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
	Success(c, http.StatusOK, gin.H{"groups": groups})
}

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
	Success(c, http.StatusOK, details)
}

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
	Success(c, http.StatusOK, gin.H{"message": "member added"})
}

func (h *Handler) deleteGroup(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}
	if err := h.services.Group.DeleteGroup(groupID, userID); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "group deleted"})
}

func (h *Handler) leaveGroup(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}
	if err := h.services.Group.LeaveGroup(groupID, userID); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "left group"})
}

func (h *Handler) confirmGroupTime(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}
	var input model.ConfirmGroupTimeDoc
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.Group.ConfirmMemberTime(groupID, userID, input.ConfirmedTime); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "time confirmed"})
}

func (h *Handler) getGroupFreeTime(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}
	var input model.GroupFreeTimeInputDoc
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "date must be in YYYY-MM-DD format")
		return
	}
	slots, err := h.services.Group.GetGroupFreeTime(groupID, userID, date)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"freeSlots": slots})
}
