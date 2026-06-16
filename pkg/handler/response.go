package handler

import (
	"druna_server/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, model.APIResponse{Data: data})
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.WithField("request_id", c.GetString(requestIDKey)).Error(message)
	c.AbortWithStatusJSON(statusCode, model.APIResponse{
		Error: &model.ErrorBody{Message: message, Code: statusCode},
	})
}
