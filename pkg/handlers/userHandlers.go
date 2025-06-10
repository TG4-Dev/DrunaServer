package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegisterHandler(c *gin.Context) { //POST
	c.JSON(http.StatusOK, gin.H{
		"message": "registering user",
	})
}

func UserLoginHandler(c *gin.Context) { // POST
	c.JSON(http.StatusOK, gin.H{
		"message": "logging in",
	})
}

func GetUserById(c *gin.Context) {
	//return
}
