package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"BlobbyServer/config"
	"BlobbyServer/pkg/handlers"
	"BlobbyServer/pkg/repositories"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/ping", handlers.PingHandler)
		api.POST("/user/register", handlers.UserRegisterHandler)
		api.POST("/user/login/requestJWT", handlers.UserLoginRequestJWTHandler)
		api.POST("/user/login/checkJWT", handlers.UserLoginCheckJWTHandler)
		api.GET("/friends", handlers.FriendListHandler)
		api.POST("/friends/request", handlers.FriendRequestHandler)
		api.POST("/events", handlers.EventAddHandler)
		api.POST("/events/free-time", handlers.EventGetFreeTime)
	}
}

func main() {
	config.LoadEnv()
	config.InitDB()

	port := config.GetEnv("PORT", "8080")

	migrate := flag.Bool("migrate", false, "Run database migrations") // go run cmd/main.go -migrate
	flag.Parse()

	if *migrate {
		log.Println("Running database migrations...")
		repositories.MigrateAll()
		log.Println("Migrations completed.")
		return
	}

	router := gin.Default()
	SetupRoutes(router)
	log.Println("Server running at localhost:" + "8080")
	log.Println("Server running at localhost:" + port)

	router.Run(":" + port)
}
