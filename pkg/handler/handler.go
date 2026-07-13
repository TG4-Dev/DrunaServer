package handler

import (
	"druna_server/pkg/service"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "druna_server/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(accessLogMiddleware())
	router.Use(metricsMiddleware())
	router.Use(corsMiddleware())

	authLimiter := NewRateLimiter(30)

	ping := router.Group("/ping")
	{
		ping.GET("/", h.ping)
	}

	if metricsEnabled() {
		registerMetricsRoute(router)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth", authLimiter.Middleware())
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/renew-token", h.renewToken)
		auth.POST("/telegram", h.telegramAuth)
	}

	apiV1 := router.Group("/api/v1", h.userIdentity)
	h.registerProtectedRoutes(apiV1)

	api := router.Group("/api", h.userIdentity)
	h.registerProtectedRoutes(api)

	return router
}

func (h *Handler) registerProtectedRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/me", h.getCurrentUser)
		users.PATCH("/me", h.updateCurrentUser)
	}

	friends := api.Group("/friends")
	{
		friends.GET("/list", h.getFriendList)
		friends.GET("/search", h.searchUsers)
		friends.GET("/request-list", h.getFriendRequestList)
		friends.GET("/requests/incoming", h.getIncomingFriendRequests)
		friends.GET("/requests/outgoing", h.getOutgoingFriendRequests)
		friends.POST("/request", h.sendFriendRequest)
		friends.POST("/accept", h.acceptFriendRequest)
		friends.POST("/reject", h.rejectFriendRequest)
		friends.DELETE("/", h.deleteFriend)
	}

	events := api.Group("/events")
	{
		events.GET("/", h.getEventList)
		events.POST("/", h.addEvent)
		events.PATCH("/:id", h.updateEvent)
		events.DELETE("/:id", h.deleteEvent)
		events.POST("/free-time", h.getFreeTime)
	}

	groups := api.Group("/groups")
	{
		groups.POST("/create", h.createGroup)
		groups.GET("/list", h.listGroups)
		groups.GET("/:id", h.getGroup)
		groups.DELETE("/:id", h.deleteGroup)
		groups.POST("/:id/leave", h.leaveGroup)
		groups.POST("/:id/members", h.addGroupMember)
		groups.POST("/:id/confirm", h.confirmGroupTime)
		groups.POST("/:id/free-time", h.getGroupFreeTime)
	}
}

func corsMiddleware() gin.HandlerFunc {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		origins = "*"
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(origins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: origins != "*",
		MaxAge:           12 * 3600,
	})
}
