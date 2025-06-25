package handler

import (
	"druna_server/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @tags Auth
// @Descrition create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body model.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 404 {object} handler.ErrorResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username     string `json:"username" binding:"required"`
	PasswordHash string `json:"passwordHash" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.PasswordHash)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
