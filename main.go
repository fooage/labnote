package main

import (
	"github.com/fooage/labnote/cache"
	"log"
	"os"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/fooage/labnote/router"
	"github.com/gin-gonic/gin"
)

// TODO: Change to a microservice architecture.

const (
	// ServerAddr is http service connection address and port.
	ServerAddr = "127.0.0.1:8090"
)

func main() {
	// Initialize and close it after calling the used database constructor.
	if err := os.Mkdir(handler.FileStorageDirectory, os.ModePerm); err != nil {
		log.Println(err)
	}
	mongodb := data.NewMongoDB()
	redis := cache.NewRedis()
	data.ConnectDatabase(mongodb)
	cache.ConnectCache(redis)
	gin.SetMode(gin.ReleaseMode)
	server := router.InitRouter(mongodb, redis)
	if err := server.Run(ServerAddr); err != nil {
		log.Println(err)
		return
	}
	data.DisconnectDatabase(mongodb)
	cache.DisconnectCache(redis)
}
