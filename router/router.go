package router

import (
	"net/http"

	"github.com/fooage/labnote/cache"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing and add local middleware. Only one main route
// is set up here to load static files and process requests, and URL of this
// website adopted consistent style.
func InitRouter(db data.Database, ch cache.Cache) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/html/*")
	router.StaticFS("/views", http.Dir("./views"))
	// set the root redirect function to the real home page
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/journal")
	})
	// set handler functions of this login page
	router.GET("/login", handler.GetLoginPage())
	router.POST("/login/submit", handler.SubmitLoginData(db))
	// set handler functions of this journal page
	router.GET("/journal", handler.VerifyAuthority(), handler.GetJournalPage())
	router.GET("/journal/list", handler.DataAuthority(db), handler.GetNotesList(db))
	router.POST("/journal/write", handler.DataAuthority(db), handler.WriteUserNote(db))
	// set handler functions of this library page
	router.GET("/library", handler.VerifyAuthority(), handler.GetLibraryPage())
	router.GET("/library/download", handler.VerifyAuthority(), handler.DownloadFile(db))
	router.GET("/library/list", handler.DataAuthority(db), handler.GetFilesList(db))
	router.GET("/library/check", handler.DataAuthority(db), handler.CheckFileStatus(ch))
	router.POST("/library/upload", handler.DataAuthority(db), handler.PostSingleChunk(ch))
	router.GET("/library/merge", handler.DataAuthority(db), handler.MergeTargetFile(db, ch))
	return router
}
