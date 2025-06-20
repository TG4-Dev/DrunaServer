package handler

import (
	"druna_server/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api", h.userIdentity)
	{
		friends := api.Group("/friends")
		{
			friends.GET("/list", h.getFriendList)
			friends.POST("/request", h.sendFriendRequest)
		}

		events := api.Group("/events")
		{
			events.POST("/list", h.getEventList)
			events.POST("/free-time", h.getFreeTime)
		}
	}
	return router
}
