package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary SignUp
// @tags Auth
// @Descrition create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body model.SignUpDoc true "account info"
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

type renewTokenInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// @Summary SignIn
// @tags Auth
// @Descrition sign in
// @ID sign in
// @Accept json
// @Produce json
// @Param input body model.SignInDoc true "account info"
// @Success 200 {integer} integer 1
// @Failure 404 {object} handler.ErrorResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.GenerateAccessRefreshToken(input.Username, input.PasswordHash)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	/*
		refreshToken, err := h.services.Authorization.GenerateRefreshToken(input.Username, input.PasswordHash)
		if err != nil {
			NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	*/
	c.JSON(http.StatusOK, map[string]interface{}{
		"access token":  accessToken,
		"refresh token": refreshToken,
	})
}

// @Summary RenewToken
// @tags Auth
// @Descrition renew token
// @ID renew token
// @Accept json
// @Produce json
// @Param input body model.RenewTokenDoc true "account info"
// @Success 200 {integer} integer 1
// @Failure 404 {object} handler.ErrorResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Failure default {object} handler.ErrorResponse
// @Router /auth/renew-token [post]
func (h *Handler) renewToken(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		NewErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userID, Username, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.RenewToken(Username, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access token":  accessToken,
		"refresh token": refreshToken,
	})
}
