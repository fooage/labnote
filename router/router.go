package router

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/proxy"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing and add local middleware.
func InitRouter(db data.Database, ch cache.Cache) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("../../views/html/*")
	router.StaticFS("../../views", http.Dir("../../views"))
	// Set the root redirect function to the real home page.
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/journal")
	})
	// These are handler functions of this login page.
	login := router.Group("/login")
	{
		login.GET("/", handler.GetLoginPage())
		login.POST("/submit", handler.SubmitLoginData(db))
	}
	// These are handler functions of this journal page.
	journal := router.Group("/journal")
	{
		journal.GET("/", handler.VerifyAuthority(), handler.GetJournalPage())
		journal.GET("/list", handler.DataAuthority(db), handler.GetNotesList(db))
		journal.POST("/write", handler.DataAuthority(db), handler.WriteUserNote(db))
	}
	// These are handler functions of this library page.
	library := router.Group("/library")
	{
		library.GET("/", handler.VerifyAuthority(), handler.GetLibraryPage())
		library.GET("/download", handler.VerifyAuthority(), handler.DownloadFile(db))
		library.GET("/list", handler.DataAuthority(db), handler.GetFilesList(db))
		library.GET("/check", handler.DataAuthority(db), handler.CheckFileStatus(ch))
		library.POST("/upload", handler.DataAuthority(db), handler.PostSingleChunk(ch))
		library.GET("/merge", handler.DataAuthority(db), handler.MergeTargetFile(db, ch))
	}
	return router
}

// Initialize the function of the reverse proxy server and implement
// multi-machine deployment through the proxy.
func InitProxy(ch cache.Cache) *gin.Engine {
	router := gin.Default()
	// Set the root redirect function to the real home page.
	router.LoadHTMLGlob("../../views/html/*")
	router.StaticFS("../../views", http.Dir("../../views"))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/journal")
	})
	// Recevie the heartbeat of servers.
	rand.Seed(time.Now().UnixNano())
	router.POST("/heartbeat", proxy.ReceiveHeartbeat())
	// These are proxy functions of this login page.
	login := router.Group("/login")
	{
		login.GET("/", proxy.UniversalReverse())
		login.POST("/submit", proxy.UniversalReverse())
	}
	// These are proxy functions of this journal page.
	journal := router.Group("/journal")
	{
		journal.GET("/", proxy.UniversalReverse())
		journal.GET("/list", proxy.UniversalReverse())
		journal.POST("/write", proxy.UniversalReverse())
	}
	// These are proxy functions of this library page.
	library := router.Group("/library")
	{
		library.GET("/", proxy.UniversalReverse())
		library.GET("/:path", proxy.FileRequestReverse(ch))
		library.POST("/:path", proxy.UploadRequestReverse(ch))
	}
	return router
}
