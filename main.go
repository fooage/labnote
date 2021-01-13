package main

import (
	"net/http"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	data.InitDatabase()
	svr := gin.Default()
	svr.LoadHTMLGlob("views/*")
	// Set the root redirect function to the real homepage.
	svr.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	// These are handler functions of this website.
	svr.GET("/home", handler.VerifyAuthority, handler.GetHomePage)
	svr.GET("/data", handler.VerifyAuthority, handler.GetNotes)
	svr.POST("/data", handler.VerifyAuthority, handler.PostNote)
	svr.GET("/login", handler.GetLoginPage)
	svr.POST("/login", handler.PostLoginData)
	svr.Run("127.0.0.1:8090")
	data.CloseDatabase()
}
