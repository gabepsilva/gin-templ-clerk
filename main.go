package main

import (
	"gotempl/middleware"
	"gotempl/views/handlers"
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
}

func main() {
	r := gin.Default()

	// Define routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/login", handlers.LoginHandler)
	r.GET("/protected", clerkMiddleware.ClerkAuthMiddleware(), handlers.ProtectedHandler)

	// Run the server
	r.Run()
}
