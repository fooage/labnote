package handler

import (
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fooage/labnote/utils"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

var (
	// CookieExpireDuration is cookie's valid duration.
	CookieExpireDuration int
	// CookieAccessScope is cookie's scope.
	CookieAccessScope string
	// FileStorageDirectory is where these files storage.
	FileStorageDirectory string
	// DownloadUrlBase decide the base url of file's url.
	DownloadUrlBase string
)

// GetJournalPage is a handler function which response the GET request for journal page.
func GetJournalPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "journal.html", gin.H{})
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

// WriteUserNote is a function that receive the log submitted in the background.
func WriteUserNote(db data.Database) gin.HandlerFunc {
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

// CheckFileStatus is a function that returns the status of the file in the server.
func CheckFileStatus(ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash
		exist, err := utils.CheckPathExists(path)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if !exist {
			// If this file not had been uploaded server will init it in the cache.
			if err := ch.ChangeFileState(hash, false); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"state": false,
					"list":  make([]cache.Chunk, 0),
				})
			}
		} else {
			saved, err := ch.CheckFileUpload(hash)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			if saved {
				// If this file had been saved in the server will return the status.
				c.JSON(http.StatusOK, gin.H{
					"state": saved,
					"list":  make([]cache.Chunk, 0),
				})
			} else {
				chunkList, err := ch.GetChunkList(hash, name)
				indexList := make([]int, 0)
				for _, chunk := range *chunkList {
					indexList = append(indexList, chunk.Index)
				}
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{})
				} else {
					c.JSON(http.StatusOK, gin.H{
						"state": saved,
						"list":  indexList,
					})
				}
			}
		}
	}
}

// PostSingleChunk is functions for receiving file slices.
func PostSingleChunk(ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.PostForm("hash")
		name := c.PostForm("name")
		path := FileStorageDirectory + "/" + hash
		blob, _ := c.FormFile("file")
		saved, err := ch.CheckFileUpload(hash)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if saved {
			c.JSON(http.StatusOK, gin.H{
				"state": saved,
				"nums":  0,
			})
		} else {
			exist, err := utils.CheckPathExists(path)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			// If there isn't a fixed path.
			if !exist {
				os.Mkdir(path, os.ModePerm)
			}
			err = c.SaveUploadedFile(blob, path+"/"+blob.Filename)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			} else {
				// If the chunk is saved successfully, insert it in the cache.
				index, _ := strconv.Atoi(blob.Filename)
				err := ch.InsertOneChunk(hash, cache.Chunk{Name: name, Hash: hash, Index: index})
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{})
					return
				}
			}
			chunkList, err := ch.GetChunkList(hash, name)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
			// Feedback the existing slices to the front.
			c.JSON(http.StatusOK, gin.H{
				"state": saved,
				"nums":  len(*chunkList),
			})
		}
	}
}

// MergeTargetFile get instructions for receiving combined files.
func MergeTargetFile(db data.Database, ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash
		exist, err := utils.CheckPathExists(path)
		if err != nil || !exist {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		// Merge the chunks to the file.
		chunkList, err := ch.GetChunkList(hash, name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		sort.Slice(*chunkList, func(a, b int) bool {
			return (*chunkList)[a].Index < (*chunkList)[b].Index
		})
		complete, _ := os.OpenFile(path+"/"+name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
		defer complete.Close()
		for _, chunk := range *chunkList {
			buffer, _ := ioutil.ReadFile(path + "/" + strconv.Itoa(chunk.Index))
			_, _ = complete.Write(buffer)
			err = os.Remove(path + "/" + strconv.Itoa(chunk.Index))
			if err != nil {
				// If an error occurs when merging files, delete the temporary files that are not fully merged.
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{})
				os.Remove(path + "/" + name)
				return
			}
		}
		// Verify file integrity.
		key, err := utils.EncodeFileHash(path, name)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		if key == hash {
			url := DownloadUrlBase + "?hash=" + hash + "&name=" + name
			if err := db.InsertOneFile(&data.File{Time: time.Now(), Name: name, Hash: hash, Url: url}); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": false,
				})
				os.RemoveAll(path)
				return
			}
			_ = ch.ChangeFileState(hash, true)
			if err = ch.RemoveAllRecords(hash); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"state": true,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"state": true,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"state": false,
			})
			os.RemoveAll(path)
		}
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

// DownloadFile is the handler function for file download.
func DownloadFile(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash + "/" + name
		exist, err := utils.CheckPathExists(path)
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
