package handlers

import (
	"github.com/gin-gonic/gin"

	"gotempl/views"
)

// Handler for the home page
func HomeHandler(c *gin.Context) {
	param := "world"
	Render(c, 200, views.Index(param))
}

// Handler for the login page
func LoginHandler(c *gin.Context) {
	Render(c, 200, views.Login())
}

// Handler for the protected page
func ProtectedHandler(c *gin.Context) {
	Render(c, 200, views.Protected())
}
