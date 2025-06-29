package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getEventList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	type id struct {
		ID int `json:"id"`
	}

	var input id
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.ID > 0 {
		userID = input.ID
	}

	eventList, err := h.services.GetEventList(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "failed to fetch event: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, eventList)
}

func (h *Handler) getFreeTime(c *gin.Context) {

}

func (h *Handler) addEvent(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
	}

	var input model.Event

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return
	}

	input.UserID = id.(int)

	eventId, err := h.services.Event.CreateEvent(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "failed to create event:"+err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"eventId": eventId,
	})

}

func (h *Handler) deleteEvent(c *gin.Context) {
	userIDInterface, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusUnauthorized, "user id not found")
		return
	}
	userID, ok := userIDInterface.(int)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id type is invalid")
		return
	}

	eventIDParam := c.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid event id")
		return
	}

	err = h.services.Event.DeleteEvent(userID, eventID)
	if err != nil {
		NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "event deleted",
	})
}
