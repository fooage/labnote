package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fooage/labnote/data"
)

var (
	// HostAddr is http service connection address and port.
	HostAddress string
	// The http server's listen port.
	ListenPort string
)

var (
	// CookieExpireDuration is cookie's valid duration.
	CookieExpireDuration int
	// CookieAccessScope is cookie's scope.
	CookieAccessScope string
	// FileStorageDirectory is where these files storage.
	FileStorageDirectory string
	// DownloadUrlBase decide the base url of file's url.
	DownloadUrlBase string
)

var (
	// TokenExpireDuration is token's valid duration.
	TokenExpireDuration time.Duration
	// EncryptionKey used for encryption.
	EncryptionKey string
	// TokenIssuer is the token's provider.
	TokenIssuer string
)

// Structure meta is my custom claims.
type meta struct {
	data.User
	jwt.StandardClaims
}

// Function generateToken is a function which generates token.
func generateToken(user data.User) (string, error) {
	claims := meta{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    TokenIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(EncryptionKey))
}

// A auxiliary function which parse token.
func parseToken(key string) (*meta, error) {
	token, err := jwt.ParseWithClaims(key, &meta{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(EncryptionKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*meta); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token not found or invalid")
}

// Function checkPathExists determine whether the path exists.
func checkPathExists(path string) (bool, error) {
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
func encodeFileHash(path string, name string) (string, error) {
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
