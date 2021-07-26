package handler

import (
	"log"
	"net/http"

	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

// GetLoginPage is a function that handles GET requests for login pages.
func GetLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

// SubmitLoginData is a function responsible for receiving verification login information.
func SubmitLoginData(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		user := &data.User{
			Email:    email,
			Password: password,
		}
		ok, err := db.CheckUserAuth(user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"pass": false, "token": nil})
			return
		}
		if ok {
			// Set the token and cookie for this user's successful login and redirect it.
			key, err := generateToken(*user)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"pass": false, "token": nil})
				return
			}
			c.SetCookie("auth", "true", CookieExpireDuration, "/", CookieAccessScope, false, true)
			c.JSON(http.StatusOK, gin.H{"pass": true, "token": key})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"pass": false, "token": nil})
	}
}
