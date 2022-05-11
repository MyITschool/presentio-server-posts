package util

import (
	"crypto/rsa"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
	"time"
)

type UserClaims struct {
	*jwt.RegisteredClaims

	ID int64
}

var accessSignKey *rsa.PrivateKey
var accessValidateKey *rsa.PublicKey

var refreshSignKey *rsa.PrivateKey
var refreshValidateKey *rsa.PublicKey

func init() {
	var err error

	accessSignKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("TOKEN_PRIVATE_KEY")))

	if err != nil {
		panic("Unable to read RSA private key")
	}

	accessValidateKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("TOKEN_PUBLIC_KEY")))

	if err != nil {
		panic("Unable to read RSA public key")
	}

	refreshSignKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")))

	if err != nil {
		panic("Unable to read RSA private key")
	}

	refreshValidateKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")))

	if err != nil {
		panic("Unable to read RSA public key")
	}
}

func CreateNewAccessToken(id int64) (string, error) {
	return createNewToken(id, &jwt.NumericDate{
		Time: time.Now().Add(time.Minute * 5),
	}, accessSignKey)
}

func CreateNewRefreshToken(id int64) (string, error) {
	return createNewToken(id, &jwt.NumericDate{
		Time: time.Now().AddDate(0, 1, 0),
	}, refreshSignKey)
}

func createNewToken(id int64, expires *jwt.NumericDate, key *rsa.PrivateKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	token.Claims = &UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: expires,
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
		ID: id,
	}

	return token.SignedString(key)
}

func ValidateAccessTokenHeader(authHeader string) (*jwt.Token, error) {
	return validateTokenHeader(authHeader, accessValidateKey)
}

func ValidateAccessToken(tokenStr string) (*jwt.Token, error) {
	return validateToken(tokenStr, accessValidateKey)
}

func ValidateRefreshTokenHeader(authHeader string) (*jwt.Token, error) {
	return validateTokenHeader(authHeader, refreshValidateKey)
}

func ValidateRefreshToken(tokenStr string) (*jwt.Token, error) {
	return validateToken(tokenStr, refreshValidateKey)
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
