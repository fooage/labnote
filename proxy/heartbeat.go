package proxy

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// The address and port of the proxy which server used!
	ProxyAddress string
	// Choose whether to use a proxy.
	UseProxy bool
	// The time limit for the registration center to determine the server timeout.
	Timeout time.Duration
)

// Send heartbeat packets, the time interval will be half shorter than the
// timeout time in order to keep the connection.
func SendHeartbeat(local string) {
	if UseProxy {
		data := url.Values{"addr": {local}}
		_, err := http.PostForm("http://"+ProxyAddress+"/heartbeat", data)
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			ticker := time.NewTicker(Timeout / 2)
			defer ticker.Stop()
			for {
				// wait a half timeout
				<-ticker.C
				data := url.Values{"addr": {local}}
				_, err := http.PostForm("http://"+ProxyAddress+"/heartbeat", data)
				if err != nil {
					return
				}
			}
		}()
	}
}

// Responsible for receiving heartbeat packets from the background server.
func ReceiveHeartbeat() gin.HandlerFunc {
	return func(c *gin.Context) {
		addr := c.PostForm("addr")
		if addr == "" {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		serverRegister(addr)
		c.JSON(http.StatusOK, gin.H{})
	}
}
