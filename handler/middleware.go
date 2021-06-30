package handler

import (
	"log"
	"net/http"

	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

// VerifyAuthority is a permission authentication middleware which verify the cookies.
func VerifyAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("auth"); err == nil {
			// find if there is a matching cookie
			if cookie == "true" {
				c.Next()
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
}

// DataAuthority function check the authentication permission for /data.
func DataAuthority(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("token")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		claims, err := parseToken(key)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			c.Abort()
			return
		}
		ok, err := db.CheckUserAuth(&claims.User)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		c.Next()
	}
}
