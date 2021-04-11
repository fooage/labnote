package data

import "time"

// User is a structure prototype for database users.
type User struct {
	Email    string `bson:"email"`    // User's email address.
	Password string `bson:"password"` // User's password.
}

// Note is a structure which stored in the database.
type Note struct {
	Time    time.Time `bson:"time"`    // Upload time of note.
	Content string    `bson:"content"` // Content of this note.
}

// Fils is a structure of server's files.
type File struct {
	Time time.Time `bson:"time"` // Upload time of file.
	Name string    `bson:"name"` // The name of this file.
	Url  string    `bson:"url"`  // Url of this file which is used to download.
}
