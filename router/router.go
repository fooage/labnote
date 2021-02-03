package router

import (
	"net/http"

	"github.com/fooage/labnote/handler"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing and add local middleware.
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/html/*")
	router.StaticFS("/views", http.Dir("./views"))
	// Set the root redirect function to the real homepage.
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/home")
	})
	// These are handler functions of this website.
	router.GET("/login", handler.GetLoginPage)
	router.POST("/login", handler.PostLoginData)
	{
		router.GET("/home", handler.VerifyAuthority(), handler.GetHomePage)
		router.GET("/data", handler.VerifyAuthority(), handler.GetNotes)
		router.POST("/data", handler.VerifyAuthority(), handler.PostNote)
	}
	return router
}
