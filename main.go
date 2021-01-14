package main

import (
	"github.com/fooage/labnote/data"
	"github.com/fooage/labnote/router"
)

// FIXME: There are some problem in redirect 301 cache.

func main() {
	data.InitDatabase()
	svr := router.InitRouter()
	svr.Run("127.0.0.1:8090")
	data.CloseDatabase()
}
