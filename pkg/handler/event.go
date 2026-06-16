package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getEventList(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	filter := model.EventFilter{Limit: 50}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}
	filter.Type = c.Query("type")
	if from := c.Query("dateFrom"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "dateFrom must be RFC3339")
			return
		}
		filter.DateFrom = &t
	}
	if to := c.Query("dateTo"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			NewErrorResponse(c, http.StatusBadRequest, "dateTo must be RFC3339")
			return
		}
		filter.DateTo = &t
	}

	result, err := h.services.Event.GetEventList(userID, filter)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "failed to fetch events: "+err.Error())
		return
	}
	Success(c, http.StatusOK, result)
}

func (h *Handler) getFreeTime(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	var input struct {
		Date string `json:"date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "date must be in YYYY-MM-DD format")
		return
	}

	slots, err := h.services.Event.GetFreeTime(userID, date)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "failed to compute free time: "+err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"freeSlots": slots})
}

func (h *Handler) addEvent(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	var input model.Event
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = userID

	eventID, err := h.services.Event.CreateEvent(input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"eventId": eventID})
}

func (h *Handler) updateEvent(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid event id")
		return
	}

	var input model.Event
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.ID = eventID
	input.UserID = userID

	if err := h.services.Event.UpdateEvent(userID, input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "event updated"})
}

func (h *Handler) deleteEvent(c *gin.Context) {
	userID := h.getUserIdFromToken(c)
	if userID == 0 {
		return
	}

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid event id")
		return
	}

	if err := h.services.Event.DeleteEvent(userID, eventID); err != nil {
		NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"message": "event deleted"})
}
