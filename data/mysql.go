package data

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Here is the definition of the database structure.
type MySQL struct {
	db   *sql.DB
	cmd  string
	name string
}

func NewMySQL(cmd string, name string) *MySQL {
	return &MySQL{
		db:   nil,
		cmd:  cmd,
		name: name,
	}
}

// InitConnection function initialize the connection to the database.
func (m *MySQL) InitConnection() error {
	conn, err := sql.Open("mysql", m.cmd+"/"+m.name)
	if err != nil {
		return err
	}
	// init variable of mysql
	m.db = conn
	err = m.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// CloseConnection is a function close the connection with mongodb.
func (m *MySQL) CloseConnection() error {
	err := m.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// CheckUserAuth is a check of if user permissions are correct.
func (m *MySQL) CheckUserAuth(user *User) (bool, error) {
	sqlStr := "select password from user where email=?"
	var password string
	err := m.db.QueryRow(sqlStr, user.Email).Scan(&password)
	if err != nil {
		// Unable to find the result, the login information is wrong.
		return false, err
	}
	if password == user.Password {
		return true, nil
	}
	return false, nil
}

// GetAllNotes function return all the notes in the database.
func (m *MySQL) GetAllNotes() (*[]Note, error) {
	sqlStr := "select time,content from note"
	rows, err := m.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var all = make([]Note, 0)
	for rows.Next() {
		var note Note
		rows.Scan(&note.Time, &note.Content)
		all = append(all, note)
	}
	return &all, nil
}

// InsertOneNote insert a new note into the database's note collection.
func (m *MySQL) InsertOneNote(note *Note) error {
	sqlStr := "insert into note(time,content) values (?,?)"
	res, err := m.db.Exec(sqlStr, note.Time, note.Content)
	if err != nil {
		return err
	}
	if res == nil {
		return err
	}
	return nil
}

// InsertOneFile insert a new file into the database's file collection.
func (m *MySQL) InsertOneFile(file *File) error {
	sqlStr := "insert into file(time,name,url) values (?,?,?)"
	res, err := m.db.Exec(sqlStr, file.Time, file.Name, file.Url)
	if err != nil {
		return err
	}
	if res == nil {
		return err
	}
	return nil
}

// GetAllFiles is used to return all of the files in storage.
func (m *MySQL) GetAllFiles() (*[]File, error) {
	sqlStr := "select time,name,url from file"
	rows, err := m.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	var all = make([]File, 0)
	for rows.Next() {
		var file File
		rows.Scan(&file.Time, &file.Name, &file.Url)
		all = append(all, file)
	}
	return &all, nil
}
