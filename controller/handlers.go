package controller

import (
	"github.com/gin-gonic/gin"

	"gotempl/views"
	"gotempl/views/layout"
)

// Handler for the home page
func HomeHandler(c *gin.Context) {
	layout.Render(c, 200, views.Index())
}

// Handler for the login page
func LoginHandler(c *gin.Context) {
	layout.Render(c, 200, views.Login())
}
