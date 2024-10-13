package layout

import (
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title   string
	Content templ.Component
	TopBar  templ.Component
	Footer  templ.Component
}

// This function will render the templ component into
// a gin context's Response Writer
func Render(c *gin.Context, status int, template templ.Component) error {

	pageData := PageData{
		Title:   "GoTempl",
		Content: template,
		TopBar:  TopBar(),
	}

	component := Layout(pageData)

	c.Status(status)
	return component.Render(c.Request.Context(), c.Writer)
}
