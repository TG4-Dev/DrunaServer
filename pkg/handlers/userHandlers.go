package handlers

import (
	"net/http"

	"BlobbyServer/pkg/services"

	"github.com/gin-gonic/gin"
)

func UserRegisterHandler(c *gin.Context) { //POST
	type RegisterInput struct {
		Name     string `json:"name" binding:"required"`
		Username string `json:"username" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := services.AuthService.Register(input.Name, input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserLoginRequestJWTHandler(c *gin.Context) { // POST
	type LoginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	token, err := services.AuthService.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func UserLoginCheckJWTHandler(c *gin.Context) { // POST
	type InputCheck struct {
		JwtToken string `json:"token" binding:"required"`
	}

	var input InputCheck
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := services.AuthService.CheckJWTService(input.JwtToken); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logging in",
	})
}

func GetUserById(c *gin.Context) {
	//return
}
