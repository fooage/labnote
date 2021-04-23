Use the web page to achieve access to the laboratory log of all platforms, and I usually want to make a log similar to a floating diary. It is now possible to store text and files that want to save.

## Deploy

Not many things are used, so it is fairly easy to deploy, mainly divided into two parts.

### Config server

There are a number of parameterized options available for customization, such as server port address, token encryption keys, and more.

```go
// main.go
const (
	// ServerAddr is http service connection address and port.
	ServerAddr = "127.0.0.1:8090"
)
// token.go
const (
	// TokenExpireDuration is token's valid duration.
	TokenExpireDuration = time.Hour * 2
	// EncryptionKey used for encryption.
	EncryptionKey = "20180212"
	// TokenIssuer is the token's provider.
	TokenIssuer = "labnote"
)
// handler.go
const (
	// CookieExpireDuration is cookie's valid duration.
	CookieExpireDuration = 7200
	// CookieAccessScope is cookie's scope.
	CookieAccessScope = "127.0.0.1"
	// FileStorageDirectory is where these files storage.
	FileStorageDirectory = "./storage"
	// DownloadUrlBase decide the base url of file's url.
	DownloadUrlBase = "http://127.0.0.1:8090/download"
)
// mongo.go
const (
	// ConnectCommand is the connection command to conncet with MongoDB.
	ConnectCommand = "mongodb://127.0.0.1:27017"
	// DatabaseName is which database will be used.
	DatabaseName = "labnote"
)
// mysql.go
const (
	//MySQL's connection statement.
	Dsn = "root:256275@tcp(127.0.0.1:3306)/labnote?charset=utf8mb4&parseTime=True"
)
```

Resumable uploads are also provided, you can change the file blocks size according to your network environment.

```javascript
// library.js
const chunkSize = 2 * 1024 * 1024;
```

### Install server

1. Use commands `go build main.go` to compile executable file.

2. Add html files in the views folder.

### Running

Registration is not available because it only used in our lab. Please enter the user directly in the database in the following format.

```json
{
  "email": "......@.....com",
  "password": "..........."
}
```

During the server running, you can use the note and file system to write or storage something. **There are some UI problem at the file system page, please do not try to break something**.

## Customized

This project is open source, welcome to customize the code and feedback questions!

### Database

In addition to providing some parameters, data interfaces are provided for use in different databases.

```go
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
	// Insert one file to this labnote system.
	InsertOneFile(file *File) error
	// Request all notes from the database.
	GetAllNotes() (*[]Note, error)
	// Request all files from the database.
	GetAllFiles() (*[]File, error)
}
```

If you want to change the database which the server used, you can change `main.go` in this snippet.

```go
db := data.NewMongoDB()
// Change database to the other.
db := data.NewMySQL()
```

### Problem

I'm not very good at web UI development, and there may be some problems here. I'm appreciate someone who asks issues and improves the code.
