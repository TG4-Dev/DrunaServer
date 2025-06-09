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
	}

	log.Println("Server running at localhost:" + "8080")

	router.Run()
}
