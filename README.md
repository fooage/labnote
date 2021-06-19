Use the web page to achieve access to the laboratory log of all platforms, and I usually want to make a log similar to a floating diary. It is now possible to store text and files that want to save.

## Deploy

Not many things are used, so it is fairly easy to deploy, mainly divided into two parts.

### Config server

There are a number of parameterized options available for customization, such as server port address, token encryption keys, and more.

```yaml
server:
  # The network address the server is running on.
  host_address: 127.0.0.1
  listen_port: 8090

  handler:
    # Cookie's valid duration the unit is seconds.
    cookie_duration: 7200
    # Cookie's determines that it can work in which domain.
    access_scope: 127.0.0.1
    # File's storage directory store files uploaded by users.
    storage_directory: ./storage
    # The default prefix when generating download links.
    url_base: http://127.0.0.1:8090/library/download

  token:
    # Token's valid duration the unit is seconds.
    token_duration: 7200
    # Key used for token encryption.
    encryption_key: 20180212
    # The issuer of the token.
    token_issuer: labnote

database:
  # Here you can choose between 'mysql' or 'mongo' options.
  type: mongo
  # Database connection command and selected database name.
  command: mongodb://admin:password@127.0.0.1:27017
  name: labnote

cache:
  address: 127.0.0.1:6379
  password: ''
  # This controls the size of the connection pool started by redis.
  pool_size: 100
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
```

If you want to change the database which the server used, you can change `config.yml`.

### Cache

In order to quickly exchange data and facilitate the storage of file slice information, we use `redis` as a cache database.

```go
// Cache interface defines the functions of the cache.
type Cache interface {
	// Init function of this cache.
	InitConnection() error
	// Close function of this cache.
	CloseConnection() error
	// Insert a chunk info in a list which key is file's hash.
	InsertOneChunk(hash string, chunk Chunk) error
	// Get the chunks list in the cache.
	GetChunkList(hash string, name string) (*[]Chunk, error)
	// Remove all of chunks after merge or error.
	RemoveAllChunks(hash string) error
	// Init the file's state for the server check.
	ChangeFileState(hash string, saved bool) error
	// Check this file whether exist completely or not.
	CheckFileUpload(hash string) (bool, error)
}
```

Like the interface provided by the database, you also can easily change the detail of functions and add functions you want.

### Problem

I'm not very good at web UI development, and there may be some problems here. I'm appreciate someone who asks issues and improves the code.
