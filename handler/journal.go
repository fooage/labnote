package handler

import (
	"html"
	"log"
	"net/http"
	"time"

	"github.com/fooage/labnote/data"
	"github.com/gin-gonic/gin"
)

// GetJournalPage is a handler function which response the GET request for journal page.
func GetJournalPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "journal.html", gin.H{})
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
