package handler

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fooage/labnote/data"
)

const (
	// TokenExpireDuration is token's valid duration.
	TokenExpireDuration = time.Hour * 2
	// EncryptionKey used for encryption.
	EncryptionKey = "20180212"
	// TokenIssuer is the token's provider.
	TokenIssuer = "labnote"
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
