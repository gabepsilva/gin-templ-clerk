package main

import (
	"gotempl/controller"
	"gotempl/database"
	"gotempl/middleware"
	"gotempl/repository"
	"gotempl/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/clerk/clerk-sdk-go/v2"

	docs "gotempl/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var clerkMiddleware middleware.ClerkPublicAuthMiddleware

func init() {
	// Load env secrets
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	apiKey := os.Getenv("CLERK_API_KEY")
	if apiKey == "" {
		panic("CLERK_API_KEY environment variable is not set")
	}

	// Load Secrets
	err = clerkMiddleware.Init()
	if err != nil {
		panic(err)

	}
	clerk.SetKey(apiKey)

	// Set logger to include the file and line number
	logrus.SetReportCaller(true)

	// Customize log format (optional)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// @title           GoTempl
// @version         1.0
// @description     My bootstrap project
// termsOfService  http://swagger.io/terms/

// contact.name   API Support
// contact.url    http://www.swagger.io/support
// contact.email  support@swagger.io

// license.name  Apache 2.0
// license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// host      localhost:8080
// BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// externalDocs.description  OpenAPI
// externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/"

	// Public routes
	// Serve static files (e.g., favicon)

	r.Static("/public", "./public")

	db := database.InitDB()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := controller.NewUserHandler(userService)

	eventRepo := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepo)
	eventHandler := controller.NewEventHandler(eventService)

	// Define routes
	eventRoutes := r.Group("/api/event")
	{
		eventRoutes.POST("/", eventHandler.CreateEvent)
		eventRoutes.GET("/", eventHandler.GetAllEvents)
		eventRoutes.GET("/:id", eventHandler.GetEvent)
		eventRoutes.PUT("/:id", eventHandler.UpdateEvent)
		eventRoutes.DELETE("/:id", eventHandler.DeleteEvent)
	}

	// User routes
	userRoutes := r.Group("/api/user")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)

	}

	adminRoutes := r.Group("/admin", clerkMiddleware.ClerkAuthMiddleware())
	{

		adminRoutes.GET("/", controller.HomeHandler)
		adminRoutes.GET("/user", userHandler.UserCRUDHandler)
		adminRoutes.GET("/event", eventHandler.EventCRUDHandler)

	}

	r.GET("/sign-in", controller.LoginHandler)

	r.GET("/swagger/*any", clerkMiddleware.ClerkAuthMiddleware(), ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run the server
	r.Run()
}
