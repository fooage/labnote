package data

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	//MySQL's connection statement.
	Dsn = "root:256275@tcp(127.0.0.1:3306)/labnote?charset=utf8mb4&parseTime=True"
)

type MySQL struct {
	db *sql.DB
}

func NewMySQL() *MySQL {
	return &MySQL{
		db: nil,
	}
}

// InitDatabase function initialize the connection to the database.
func (m *MySQL) InitDatabase() error {
	conn, err := sql.Open("mysql", Dsn)
	if err != nil {
		return err
	}
	// Init these variable of mongodb.
	m.db = conn
	err = m.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// CloseDatabase is a function close the connection with mongodb.
func (m *MySQL) CloseDatabase() error {
	err := m.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// CheckUserAuth is a check of if user permissions are correct.
func (m *MySQL) CheckUserAuth(user *User) (bool, error) {
	sqlStr := "select password from user where email=?"
	var u string
	err := m.db.QueryRow(sqlStr, user.Email).Scan(&u)
	if err != nil {
		// Unable to find the result, the login information is wrong.
		return false, err
	}
	if u == user.Password {
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
		var n Note
		rows.Scan(&n.Time, &n.Content)
		all = append(all, n)
	}
	return &all, nil
}

// InsertOneNote insert a new note into the database's note collection.
func (m *MySQL) InsertOneNote(note *Note) error {
	sqlStr := "insert into note(time,content) values (?,?)"
	ret, err := m.db.Exec(sqlStr, note.Time, note.Content)
	if err != nil {
		return err
	}
	if ret == nil {
		return err
	}
	return nil
}

// InsertOneFile insert a new file into the database's file collection.
func (m *MySQL) InsertOneFile(file *File) error {
	sqlStr := "insert into file(time,name,url) values (?,?,?)"
	ret, err := m.db.Exec(sqlStr, file.Time, file.Name, file.Url)
	if err != nil {
		return err
	}
	if ret == nil {
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
		var f File
		rows.Scan(&f.Time, &f.Name, &f.Url)
		all = append(all, f)
	}
	return &all, nil
}
