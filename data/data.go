package data

// Database definition of the abstract interface of it.
type Database interface {
	// The init function of this database.
	InitConnection() error
	// The close function of this database.
	CloseConnection() error
	// Verify that the user information is reasonable.
	CheckUserAuth(user *User) (bool, error)
	// Insert one note to this labnote system.
	InsertOneNote(note *Note) error
	// Insert one file to this labnote system.
	InsertOneFile(file *File) error
	// Request all notes from the database.
	GetAllNotes() (*[]Note, error)
	// Request all files from the database.
	GetAllFiles() (*[]File, error)
}

// ConnectDatabase is function which load the database.
func ConnectDatabase(data Database) error {
	err := data.InitConnection()
	if err != nil {
		return err
	}
	return nil
}

// DisconnectDatabase is function which disconnect the database.
func DisconnectDatabase(data Database) error {
	err := data.CloseConnection()
	if err != nil {
		return err
	}
	return nil
}
