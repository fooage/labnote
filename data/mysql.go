package data

// TODO: Prepare to supplement the MySQL database code and switch the background database.

// Here is the definition of the database structure.
type MySQL struct{}

func NewMySQL() *MySQL {
	return &MySQL{}
}

// InitDatabase function initialize the connection to the database.
func (m *MySQL) InitDatabase() error { return nil }

// CloseDatabase is a function close the connection with mongodb.
func (m *MySQL) CloseDatabase() error { return nil }

// CheckUserAuth is a check of if user permissions are correct.
func (m *MySQL) CheckUserAuth(user *User) (bool, error) { return true, nil }

// InsertOneNote insert a new note into the database's note collection.
func (m *MySQL) InsertOneNote(note *Note) error { return nil }

// GetAllNotes function return all the notes in the database.
func (m *MySQL) GetAllNotes() (*[]Note, error) { return nil, nil }
