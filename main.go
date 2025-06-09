package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"BlobbyServer/handlers"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/ping", handlers.PingHandler)

		api.POST("/user/register", handlers.UserRegisterHandler)

		api.POST("/user/login", handlers.UserLoginHandler)

		api.GET("/friends", handlers.FriendListHandler)

		api.POST("/friends/request", handlers.FriendRequestHandler)

		api.POST("/events", handlers.EventAddHandler)

		api.POST("/events/free-time", handlers.EventGetFreeTime)
	}

	log.Println("Server running at localhost:" + "8080")

	router.Run()
}
