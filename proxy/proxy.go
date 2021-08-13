package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/fooage/labnote/cache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	// HostAddr is proxy service connection address and port!
	HostAddress string
	// The proxy server's listen port!
	ListenPort string
	// Choose the http's version of proxy.
	HttpVersion string
)

// UniversalReverse random reverse proxy for ordinary requests to complete simple reverse.
func UniversalReverse() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get a random target address
		target := getRandomServer()
		if target == "" {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		director := func(req *http.Request) {
			req.URL.Scheme = HttpVersion
			req.URL.Host = target
			req.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// A more detailed proxy is needed for requests with file hash values, such as
// redirection based on the IP of the server where the file is stored.
func FileRequestReverse(ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		var target string
		switch c.Param("path") {
		case "list":
			target = getRandomServer()
		case "check":
			location, err := ch.GetFileLocation(c.Query("hash"))
			if err == redis.Nil {
				target = getRandomServer()
			} else if err == nil {
				target = location
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		case "download", "merge":
			// TODO: To code the p2p file system and change "download" action
			//  to the random. It is very difficult to complete a distributed
			//  file system, but I will try to do it.
			location, err := ch.GetFileLocation(c.Query("hash"))
			if err == redis.Nil {
				c.JSON(http.StatusBadRequest, gin.H{})
				return
			} else if err == nil {
				target = location
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}
		if target != "" {
			// Target Address is not nil, reverse the requeset.
			director := func(request *http.Request) {
				request.URL.Scheme = HttpVersion
				request.URL.Host = target
				request.Host = target
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(c.Writer, c.Request)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
		}
	}
}

func UploadRequestReverse(ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		var target string
		location, err := ch.GetFileLocation(c.Query("hash"))
		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		} else if err == nil {
			target = location
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if target != "" {
			// Target Address is not nil, reverse the requeset.
			director := func(request *http.Request) {
				request.URL.Scheme = HttpVersion
				request.URL.Host = target
				request.Host = target
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(c.Writer, c.Request)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
		}
	}
}
