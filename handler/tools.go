package handler

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

// FileHash return file md5 code generation function is used to verify integrity.
func FileHash(path string, name string) (string, error) {
	file, err := os.Open(path + "/" + name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	io.Copy(hash, file)
	key := hex.EncodeToString(hash.Sum(nil))
	return key, nil
}

// Structure definition of file slice.
type Chunk struct {
	Index int         // Index of chunk.
	Info  fs.FileInfo // Chunk's info.
}

// MergeSlice is function merge the file slice.
func MergeSlice(path string, name string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	chunks := make([]Chunk, 0)
	for _, info := range files {
		idx, _ := strconv.Atoi(info.Name())
		chunks = append(chunks, Chunk{
			Index: idx,
			Info:  info,
		})
	}
	sort.Slice(chunks, func(a, b int) bool {
		return chunks[a].Index < chunks[b].Index
	})
	complete, _ := os.OpenFile(path+"/"+name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer complete.Close()
	for _, chunk := range chunks {
		buffer, err := ioutil.ReadFile(path + "/" + chunk.Info.Name())
		if err != nil {
			return err
		}
		complete.Write(buffer)
		os.Remove(path + "/" + chunk.Info.Name())
	}
	return nil
}

// Determine whether the path exists.
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
