package main

import (
	"log"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/config"
	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/proxy"
	"github.com/fooage/labnote/router"
	"github.com/gin-gonic/gin"
)

func main() {
	db, ch, err := config.ProxyConfiguration()
	if err != nil {
		log.Println(err)
		return
	} else {
		// Initialize and close it after calling the used database constructor.
		data.ConnectDatabase(db)
		defer data.DisconnectDatabase(db)
		cache.ConnectCache(ch)
		defer cache.DisconnectCache(ch)
	}
	gin.SetMode(gin.DebugMode)
	server := router.InitProxy(ch)
	// Begin to proxy the request from users.
	if err := server.Run(proxy.HostAddress + ":" + proxy.ListenPort); err != nil {
		log.Println(err)
		return
	}
}
