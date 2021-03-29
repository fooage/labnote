package data

import "fmt"

// The definition of the abstract interface of the database.
type Database interface {
	// The init function of this database.
	InitDatabase() error
	// The close function of this database.
	CloseDatabase() error
	// Verify that the user information is reasonable.
	CheckUserAuth(user *User) (bool, error)
	// Insert one note to this labnote system.
	InsertOneNote(note *Note) error
	// Request all notes from the database.
	GetAllNotes() (*[]Note, error)
}

// ConnectDatabase is function which load the database.
func ConnectDatabase(data Database) {
	err := data.InitDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// DisconnectDatabase is function which disconnect the database.
func DisconnectDatabase(data Database) {
	err := data.CloseDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
}
