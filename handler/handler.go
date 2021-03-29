package handler

import (
	"html"
	"net/http"
	"time"

	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

// CookieExpireDuration is cookie's valid duration.
var CookieExpireDuration = 7200

// VerifyAuthority is a permission authentication middleware.
func VerifyAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("auth"); err == nil {
			// Find if there is a matching cookie here.
			if cookie == "true" {
				c.Next()
				return
			}
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
}

// GetHomePage is a handler function which response the GET request for homepage.
func GetHomePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	}
}

// GetLoginPage is a function that handles GET requests for login pages.
func GetLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}

// PostLoginData is a function responsible for receiving verification login information.
func PostLoginData(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		user := &data.User{
			Email:    email,
			Password: password,
		}
		res, _ := db.CheckUserAuth(user)
		if res {
			// Set the token and cookie for this user's successful login and redirect it.
			key, _ := GenerateToken(*user)
			c.SetCookie("auth", "true", CookieExpireDuration, "/", "127.0.0.1", false, true)
			c.JSON(http.StatusOK, gin.H{"pass": true, "token": key})
		} else {
			c.JSON(http.StatusOK, gin.H{"pass": false, "token": nil})
		}
	}
}

// GetNotes get all the notes it have so far.
func GetNotes(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		notes, err := db.GetAllNotes()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		c.JSON(http.StatusOK, gin.H{"notes": *notes})
	}
}

// DataAuthority function check the authentication permission for /data.
func DataAuthority(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("token")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{})
			c.Abort()
			return
		}
		claims, err := ParseToken(key)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			c.Abort()
			return
		}
		res, _ := db.CheckUserAuth(&claims.User)
		if !res {
			c.JSON(http.StatusBadRequest, gin.H{})
			c.Abort()
			return
		}
		c.Next()
	}
}

// PostNote is a function that receive the log submitted in the background.
func PostNote(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		content := c.PostForm("content")
		note := &data.Note{
			Time:    time.Now(),
			Content: html.EscapeString(content),
			// Escaping to prevent XSS attacks.
		}
		err := db.InsertOneNote(note)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
