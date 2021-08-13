package main

import (
	"log"
	"os"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/config"
	"github.com/fooage/labnote/proxy"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/fooage/labnote/router"
	"github.com/gin-gonic/gin"
)

// TODO: Change to a microservice architecture.

func main() {
	db, ch, err := config.ServerConfiguration()
	if err != nil {
		log.Println(err)
		return
	} else {
		// Initialize and close it after calling the used database constructor.
		if err := os.Mkdir(handler.FileStorageDirectory, os.ModePerm); err != nil {
			log.Println(err)
		}
		data.ConnectDatabase(db)
		defer data.DisconnectDatabase(db)
		cache.ConnectCache(ch)
		defer cache.DisconnectCache(ch)
	}
	// Start sending requests to proxy server, if the proxy server option is
	// checked. If there is no proxy server, server can be run independently.
	proxy.SendHeartbeat(handler.HostAddress + ":" + handler.ListenPort)
	gin.SetMode(gin.DebugMode)
	server := router.InitRouter(db, ch)
	if err := server.Run(handler.HostAddress + ":" + handler.ListenPort); err != nil {
		log.Println(err)
		return
	}
}
