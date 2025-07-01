package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get event List
// @Security ApiKeyAuth
// @Tags events
// @Description Get current user's event list
// @ID get-events
// @Produce  json
// @Success 200 {array} model.EventDoc
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/events [get]
func (h *Handler) getEventList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	type id struct {
		ID int `json:"id"`
	}

	var input id

	idStr := c.Query("id")
	if idStr != "" { // Если ничего не передано в запросе, вернёт список по id из token
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "ID must be a number"})
			return
		}
		input.ID = id
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

// @Summary Create Event
// @Security ApiKeyAuth
// @Tags events
// @Description Create event
// @ID create-event
// @Accept  json
// @Produce  json
// @Param input body model.EventDoc true "list info"
// @Success 200 {object} model.AddEventDoc
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/events/ [post]
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

// @Summary Delete Event
// @Security ApiKeyAuth
// @Tags events
// @Description Create event
// @ID delete-event
// @Accept  json
// @Produce  json
// @Param input body model.DeleteEventDoc true "list info"
// @Success 200 {object} model.AddEventDoc
// @Failure 400,404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /api/events/ [delete]
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
