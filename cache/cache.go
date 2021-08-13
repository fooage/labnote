package cache

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

// ConnectCache is function which load the database.
func ConnectCache(cache Cache) error {
	err := cache.InitConnection()
	if err != nil {
		return err
	}
	return nil
}

// DisconnectCache is function which disconnect the database.
func DisconnectCache(cache Cache) error {
	err := cache.CloseConnection()
	if err != nil {
		return err
	}
	return nil
}
