package cache

// Chunk is the definition of a part of file.
// Remember the 'Hash' in this structure from the file not this hash.
type Chunk struct {
	Name  string // name of the file this chunk belongs to.
	Hash  string // file hash to which this chunk belongs.
	Index int    // index of chunk in the file.
}

type Target struct {
	Hash     string // the hash of this file.
	Location string // the storage server's host.
}
