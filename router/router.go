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
	// Set the root redirect function to the real home page.
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/journal")
	})
	// These are handler functions of this website.
	router.GET("/login", handler.GetLoginPage())
	router.POST("/login/submit", handler.SubmitLoginData(db))

	router.GET("/journal", handler.VerifyAuthority(), handler.GetJournalPage())
	router.GET("/journal/list", handler.DataAuthority(db), handler.GetNotesList(db))
	router.POST("/journal/write", handler.DataAuthority(db), handler.WriteUserNote(db))

	router.GET("/library", handler.VerifyAuthority(), handler.GetLibraryPage())
	router.GET("/library/download", handler.VerifyAuthority(), handler.DownloadFile(db))
	router.GET("/library/list", handler.DataAuthority(db), handler.GetFilesList(db))
	router.GET("/library/check", handler.DataAuthority(db), handler.CheckFileStatus())
	router.POST("/library/upload", handler.DataAuthority(db), handler.PostSingleChunk(db))
	router.GET("/library/merge", handler.DataAuthority(db), handler.MergeTargetFile(db))

	return router
}
