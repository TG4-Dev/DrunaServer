package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @Security ApiKeyAuth
// @tags events
// @Descrition get event list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body model.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400, 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) getEventList(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	userID, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return
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
		return
	}

	var input model.Event

	err := c.BindJSON(&input)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := id.(int)
	if !ok {
		NewErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return
	}

	input.UserID = strconv.Itoa(userID)

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
