package handler

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

const (
	// CookieExpireDuration is cookie's valid duration.
	CookieExpireDuration = 7200
	// CookieAccessScope is cookie's scope.
	CookieAccessScope = "127.0.0.1"
	// FileStorageDirectory is where these files storage.
	FileStorageDirectory = "./storage"
	// DownloadUrlBase decide the base url of file's url.
	DownloadUrlBase = "http://127.0.0.1:8090/download"
)

// VerifyAuthority is a permission authentication middleware.
func VerifyAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("auth"); err == nil {
			log.Println(err)
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

// GetLibraryPage is the function which response to the library page's get request.
func GetLibraryPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "library.html", gin.H{})
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
		ok, err := db.CheckUserAuth(user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"pass": false, "token": nil})
			return
		}
		if ok {
			// Set the token and cookie for this user's successful login and redirect it.
			key, err := GenerateToken(*user)
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

// GetNotesList get all the notes it have so far.
func GetNotesList(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		notes, err := db.GetAllNotes()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"notes": nil})
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
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
			return
		}
		claims, err := ParseToken(key)
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
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}

// GetChunkList is a function that returns the status of the file in the server.
func GetChunkList(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash
		chunkList := []string{}
		exist, err := PathExists(path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		state := false
		if exist {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			for _, file := range files {
				fileName := file.Name()
				chunkList = append(chunkList, fileName)
				if fileName == name {
					state = true
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"state": state,
			"list":  chunkList,
		})
	}
}

// PostChunk is functions for receiving file slices.
func PostChunk(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.PostForm("hash")
		name := c.PostForm("name")
		path := FileStorageDirectory + "/" + hash
		chunk, _ := c.FormFile("file")
		exist, err := PathExists(path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		// If there isn't a fixed path.
		if !exist {
			os.Mkdir(path, os.ModePerm)
		}
		err = c.SaveUploadedFile(chunk, FileStorageDirectory+"/"+hash+"/"+chunk.Filename)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		state := false
		chunkList := []string{}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		for _, file := range files {
			fileName := file.Name()
			chunkList = append(chunkList, fileName)
			if fileName == name {
				state = true
			}
		}
		// Feedback the existing slices to the front.
		c.JSON(http.StatusOK, gin.H{
			"state": state,
			"list":  chunkList,
		})
	}
}

// GetMergeFile get instructions for receiving combined files.
func GetMergeFile(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash
		exist, err := PathExists(path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if !exist {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if err := MergeSlice(path, name); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		// Verify file integrity.
		key, err := FileHash(path, name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if key == hash {
			url := DownloadUrlBase + "?hash=" + hash + "&name=" + name
			if err := db.InsertOneFile(&data.File{
				Time: time.Now(),
				Name: name,
				Url:  url,
			}); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": false,
				})
				os.RemoveAll(path)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"state": true,
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"state": false,
		})
		os.RemoveAll(path)
	}
}

// GetFilesList get all the files in server's storage.
func GetFilesList(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		files, err := db.GetAllFiles()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, gin.H{"files": files})
	}
}

// GetFile is the handler function for file download.
func GetFile(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash + "/" + name
		exist, err := PathExists(path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if !exist {
			c.JSON(http.StatusInternalServerError, gin.H{})
		} else {
			c.File(path)
		}
	}
}
