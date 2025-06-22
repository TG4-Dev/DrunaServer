package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getEventList(c *gin.Context) {

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
}
