package router

import (
	"net/http"

	"github.com/fooage/labnote/handler"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing and add local middleware.
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	// Set the root redirect function to the real homepage.
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	// These are handler functions of this website.
	router.GET("/login", handler.GetLoginPage)
	router.POST("/login", handler.PostLoginData)
	router.Use(handler.VerifyAuthority())
	{
		router.GET("/home", handler.GetHomePage)
		router.GET("/data", handler.GetNotes)
		router.POST("/data", handler.PostNote)
	}
	return router
}
