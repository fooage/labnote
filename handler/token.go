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

// Claims is my custom claims.
type Claims struct {
	data.User
	jwt.StandardClaims
}

// GenerateToken is a function which generates token.
func GenerateToken(user data.User) (string, error) {
	claims := Claims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    TokenIssuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(EncryptionKey))
}

// ParseToken is a function which parse token.
func ParseToken(key string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(key, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(EncryptionKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token not found or invalid")
}
