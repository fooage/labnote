package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fooage/labnote/cache"
	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

// GetLibraryPage is the function which response to the library page's get request.
func GetLibraryPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "library.html", gin.H{})
	}
}

// CheckFileStatus is a function that returns the status of the file in the server.
func CheckFileStatus(ch cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash
		exist, err := checkPathExists(path)
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
			// init file' location
			ch.InitFileLocation(hash, HostAddress+":"+ListenPort)
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
			exist, err := checkPathExists(path)
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
		exist, err := checkPathExists(path)
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
		key, err := encodeFileHash(path, name)
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

// DownloadFile is the handler function for file download.
func DownloadFile(db data.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Query("hash")
		name := c.Query("name")
		path := FileStorageDirectory + "/" + hash + "/" + name
		exist, err := checkPathExists(path)
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
