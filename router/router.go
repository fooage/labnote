package router

import (
	"net/http"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing and add local middleware.
func InitRouter(db data.Database) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/html/*")
	router.StaticFS("/views", http.Dir("./views"))
	// Set the root redirect function to the real homepage.
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/home")
	})
	// These are handler functions of this website.
	router.GET("/login", handler.GetLoginPage())
	router.POST("/login", handler.PostLoginData(db))
	{
		router.GET("/home", handler.VerifyAuthority(), handler.GetHomePage())
		router.GET("/library", handler.VerifyAuthority(), handler.GetLibraryPage())
		router.GET("/download", handler.VerifyAuthority(), handler.GetFile(db))
	}
	{
		router.GET("/note", handler.DataAuthority(db), handler.GetNotesList(db))
		router.POST("/write", handler.DataAuthority(db), handler.PostNote(db))
		router.GET("/file", handler.DataAuthority(db), handler.GetFilesList(db))
		router.GET("/check", handler.DataAuthority(db), handler.GetChunkList(db))
		router.POST("/upload", handler.DataAuthority(db), handler.PostChunk(db))
		router.GET("/merge", handler.DataAuthority(db), handler.GetMergeFile(db))
	}
	return router
}
