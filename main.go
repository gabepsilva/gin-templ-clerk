package main

import (
	"gotempl/database"
	"gotempl/handlers"
	"gotempl/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/clerk/clerk-sdk-go/v2"
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

	database.InitDB()

}

func main() {
	r := gin.Default()

	eventHandler := handlers.NewEventHandler()
	userHandler := handlers.NewUserHandler()

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

		adminRoutes.GET("/", handlers.HomeHandler)
		adminRoutes.GET("/protected", handlers.ProtectedHandler)
		adminRoutes.GET("/user", handlers.NewUserHandler().UserCRUDHandler)

	}

	r.GET("/sign-in", handlers.LoginHandler)

	// Run the server
	r.Run()
}
