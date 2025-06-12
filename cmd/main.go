package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"

	"BlobbyServer/pkg/handlers"
	"BlobbyServer/pkg/repositories"
)

func SetupRoutes(router *gin.Engine) {
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
}

func main() {

	migrate := flag.Bool("migrate", false, "Run database migrations") // go run cmd/main.go -migrate
	flag.Parse()

	if *migrate {
		log.Println("Running database migrations...")
		repositories.MigrateAll()
		log.Println("Migrations completed.")
		return
	}

	// db, err := storage.NewConnection() // idk where and how to connect to the db

	// if err != nil {
	// 	log.Fatal("could not loat the database")
	// }

	router := gin.Default()
	SetupRoutes(router)
	log.Println("Server running at localhost:" + "8080")

	router.Run()
}
