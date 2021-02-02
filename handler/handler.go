package handler

import (
	"html"
	"net/http"
	"time"

	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

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
func GetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

// GetLoginPage is a function that handles GET requests for login pages.
func GetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// PostLoginData is a function responsible for receiving verification login information.
func PostLoginData(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	user := &data.User{
		Email:    email,
		Password: password,
	}
	if data.CheckUserAuth(user) {
		// Set the cookie for this user's successful login and redirect it.
		c.SetCookie("auth", "true", 3600, "/", "127.0.0.1", false, true)
		c.JSON(http.StatusOK, gin.H{"pass": true, "email": email, "password": password})
	} else {
		c.JSON(http.StatusOK, gin.H{"pass": false, "email": email, "password": password})
	}
}

// GetNotes get all the notes it have so far.
func GetNotes(c *gin.Context) {
	notes, err := data.GetAllNotes()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notes": *notes})
}

// PostNote is a function that receive the log submitted in the background.
func PostNote(c *gin.Context) {
	content := c.PostForm("content")
	note := &data.Note{
		Time:    time.Now(),
		Content: html.EscapeString(content),
		// Escaping to prevent XSS attacks.
	}
	err := data.InsertOneNote(note)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
