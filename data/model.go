package data

import "time"

// User is a structure prototype for database users.
type User struct {
	Email    string `bson:"email" json:"email"`       // User's email address.
	Password string `bson:"password" json:"password"` // User's password.
}

// Note is a structure which stored in the database.
type Note struct {
	Time    time.Time `bson:"time"`    // Upload time of note.
	Content string    `bson:"content"` // Content of this note.
}
