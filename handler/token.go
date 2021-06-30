package handler

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fooage/labnote/data"
)

var (
	// token's valid duration
	TokenExpireDuration time.Duration
	// key used for encryption
	EncryptionKey string
	// token's provider
	TokenIssuer string
)

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

// A auxiliary function which parse token to get the claims information.
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
