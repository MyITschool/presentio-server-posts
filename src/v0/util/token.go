package util

import (
	"crypto/rsa"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
)

type UserClaims struct {
	*jwt.RegisteredClaims

	ID int64
}

var accessValidateKey *rsa.PublicKey

func init() {
	var err error

	accessValidateKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("TOKEN_PUBLIC_KEY")))

	if err != nil {
		panic("Unable to read RSA public key")
	}
}

func ValidateAccessTokenHeader(authHeader string) (*jwt.Token, error) {
	return validateTokenHeader(authHeader, accessValidateKey)
}

func validateTokenHeader(authHeader string, key *rsa.PublicKey) (*jwt.Token, error) {
	if authHeader == "" {
		return nil, jwt.ErrTokenMalformed
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, jwt.ErrTokenMalformed
	}

	return validateToken(parts[1], key)
}

func validateToken(tokenStr string, key *rsa.PublicKey) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}, jwt.WithValidMethods([]string{"RS256"}))
}

func HandleTokenError(err error, c *gin.Context) {
	if errors.Is(err, jwt.ErrTokenMalformed) {
		c.Status(406)
	} else if errors.Is(err, jwt.ErrTokenExpired) {
		c.Status(408)
	} else {
		c.Status(400)
	}
}
