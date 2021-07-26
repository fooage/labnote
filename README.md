Use the web page to achieve access to the laboratory log of all platforms, and I usually want to make a log similar to a diary. It is now possible to store text and files that want to save. **Now you can deploy a file document storage service in a cluster.**

## Install server

The system is mainly divided into two parts, reverse proxy service and main service, in these two folders `/cmd/proxy` and `/cmd/server`. The main service can run independently without setting up a proxy!

1. `git clone` get all of these source code.

2. Use commands `go build main.go` to compile executable file.

## Configuration

When deploying the service, almost all settings are in the configuration file, which can be quickly changed and customized.

### Config server

```yaml
server:
  # If you use a proxy for multi-machine deployment, you need to change the
  # address settings and port settings, and set the proxy server address and
  # select an appropriate timeout duration.
  handler:
    # The network address the server is running on.
    host_address: 127.0.0.1
    listen_port: 8090
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

  proxy:
    proxy_address: 127.0.0.1:11000
    # Choose whether use the proxy server.
    use_proxy: true
    # Timeout duration of the registry.
    timeout: 120

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

### Proxy Config

```yaml
proxy:
  # When using a proxy, static files will be loaded by the proxy server, and be
  # sure to select a suitable timeout period to ensure the stability of the system.
  universal:
    # The network address the proxy is running on.
    host_address: 127.0.0.1
    listen_port: 11000
    # The http version of this proxy such as "http" and "https".
    http_version: 'http'
  registry:
    # Timeout duration of the registry.
    timeout: 120

database:
  # Here you can choose between 'mysql' or 'mongo' options.
  type: mongo
  # Database connection command and selected database name.
  command: mongodb://admin:password@127.0.0.1:27017
  name: labnote

cache:
  # Cache the authentication information of the database.
  address: 127.0.0.1:6379
  password: ''
  # This controls the size of the connection pool started by redis.
  pool_size: 100
```

### Transform Config

```javascript
// Resumable uploads are also provided, you can change the file blocks size
// according to your network environment.
const chunkSize = 2 * 1024 * 1024;
```

## Running

**If you need to use a proxy server, you need to start it before starting the main service.** And registration is not available because it only used in our lab or small-scale use. So please enter the user directly in the database in the following format.

```json
{
  "email": "......@.....com",
  "password": "..........."
}
```

During the server running, you can use the note and file system to write or storage something. ~~There are some UI problem at the file system page, please do not try to break something.~~

## Customized

This project is open source, welcome to customize the code and feedback questions! I try to make the system modular as much as possible to simplify the addition of later modules.

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
	RemoveAllRecords(hash string) error
	// Init the file's state for the server check.
	ChangeFileState(hash string, saved bool) error
	// Check this file whether exist completely or not.
	CheckFileUpload(hash string) (bool, error)
	// Initialize the location of this file for reverse.
	InitFileLocation(hash string, addr string) error
	// Get the storage location of the file and get where it is.
	GetFileLocation(hash string) (string, error)
}
```

Like the interface provided by the database, you also can easily change the detail of functions and add functions you want.

## Problem

I'm not very good at web UI development, and there may be some problems here. I'm appreciate someone who asks issues and improves the code. And I have an idea to write a distributed file system based on p2p network, but don't know how to do it.
