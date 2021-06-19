package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// CheckPathExists determine whether the path exists.
func CheckPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// EncodeFileHash return file md5 code generation function is used to verify integrity.
func EncodeFileHash(path string, name string) (string, error) {
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
