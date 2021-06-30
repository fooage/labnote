package main

import (
	"log"
	"os"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/config"
	"github.com/spf13/viper"

	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/handler"
	"github.com/fooage/labnote/router"
	"github.com/gin-gonic/gin"
)

// TODO: Change to a microservice architecture.

var (
	// http service connection address
	HostAddress string
	// http server's listen port
	ListenPort string
)

// Read the config file and initialize server address meta information.
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
		return
	}
	HostAddress = viper.GetString("server.host_address")
	ListenPort = viper.GetString("server.listen_port")
}

func main() {
	db, ch, err := config.LoadConfiguration()
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
	// set the mode of engine
	gin.SetMode(gin.ReleaseMode)
	server := router.InitRouter(db, ch)
	if err := server.Run(HostAddress + ":" + ListenPort); err != nil {
		log.Println(err)
		return
	}
}
