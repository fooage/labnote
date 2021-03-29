package main

import (
	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/router"
)

// TODO: It will add the support of others database.

func main() {
	// Initialize and close it after calling the used database constructor.
	db := data.NewMongoDB()
	data.ConnectDatabase(db)
	svr := router.InitRouter(db)
	svr.Run(":8090")
	data.DisconnectDatabase(db)
}
