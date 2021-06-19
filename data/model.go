package data

import "time"

// User is a structure prototype for database users.
type User struct {
	Email    string `bson:"email"`    // user's email address.
	Password string `bson:"password"` // user's password.
}

// Note is a structure which stored in the database.
type Note struct {
	Time    time.Time `bson:"time"`    // upload time of note.
	Content string    `bson:"content"` // content of this note.
}

// File is a structure of server's files.
type File struct {
	Time time.Time `bson:"time"` // upload time of file.
	Name string    `bson:"name"` // the name of this file.
	Hash string    `bson:"hash"` // the hash of file saved.
	Url  string    `bson:"url"`  // url of this file which is used to download.
}
