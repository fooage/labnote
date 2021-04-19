package main

import (
	"os"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/fooage/labnote/router"
	"github.com/gin-gonic/gin"
)

// TODO: It will add the support of others database.

const (
	// ServerAddr is http service connection address and port.
	ServerAddr = "127.0.0.1:8090"
)

func main() {
	// Initialize and close it after calling the used database constructor.
	os.Mkdir(handler.FileStorageDirectory, os.ModePerm)
	db := data.NewMongoDB()
	data.ConnectDatabase(db)
	gin.SetMode(gin.ReleaseMode)
	svr := router.InitRouter(db)
	svr.Run(ServerAddr)
	data.DisconnectDatabase(db)
}
