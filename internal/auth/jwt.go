package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

func GenerateJwt(userId int, expiresInSeconds *int, jwtSecret string) (string, error) {
	expirationTime := 3600 // default expiration time in seconds (1 hour)
	if expiresInSeconds != nil && *expiresInSeconds < expirationTime {
		expirationTime = *expiresInSeconds
	}

	currentTime := time.Now()
	expiredAt := currentTime.Add(time.Second * time.Duration(expirationTime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(currentTime),
		ExpiresAt: jwt.NewNumericDate(expiredAt),
		Subject:   fmt.Sprint(userId),
	}})
	ss, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseJwt(tokenString, jwtSecret string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}

func GetApiToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitToken := strings.Split(authHeader, "ApiKey ")
	if len(splitToken) != 2 {
		return "", errors.New("malformed authorization header")
	}
	return splitToken[1], nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("malformed authorization header")
	}
	return splitToken[1], nil
}
