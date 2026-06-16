package handler

import (
	"druna_server/pkg/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.PasswordHash == "" && input.Password != "" {
		input.PasswordHash = input.Password
	}
	if input.PasswordHash == "" && input.Password == "" {
		NewErrorResponse(c, http.StatusBadRequest, "password is required")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, http.StatusOK, gin.H{"id": id})
}

type signInInput struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
}

type renewTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

func (in signInInput) password() string {
	if in.Password != "" {
		return in.Password
	}
	return in.PasswordHash
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if input.password() == "" {
		NewErrorResponse(c, http.StatusBadRequest, "password or passwordHash is required")
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.GenerateAccessRefreshToken(input.Username, input.password())
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	Success(c, http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) renewToken(c *gin.Context) {
	var input renewTokenInput
	_ = c.ShouldBindJSON(&input)

	token := input.RefreshToken
	if token == "" {
		header := c.GetHeader(authorizationHeader)
		if header == "" {
			NewErrorResponse(c, http.StatusUnauthorized, "empty auth header or refreshToken in body")
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
			return
		}
		token = headerParts[1]
	}
	if token == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.RenewToken(token)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	Success(c, http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) telegramAuth(c *gin.Context) {
	var input struct {
		InitData string `json:"initData" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.LoginWithTelegramInitData(input.InitData)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "telegram auth failed: "+err.Error())
		return
	}

	Success(c, http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
