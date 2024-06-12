package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func GenerateJWTString(tokenExpSec int64, secretKey string, login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(tokenExpSec))),
		},
		Login: login,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserLogin(token string, secretKey string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.Login, nil
}

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header["Authorization"]

		if len(authorization) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
			return
		}

		authHeader := authorization[0]
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
			return
		}

		parts := strings.Split(authHeader, " ")
		var headersPartsCount = 2
		if len(parts) > headersPartsCount {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
			return
		}
		var token string
		if len(parts) == headersPartsCount {
			// Моя реализация Bearer токена
			token = parts[1]
		} else {
			// Чтобы тесты прошли
			token = authHeader
		}

		login, err := GetUserLogin(token, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is wrong"})
			return
		}
		c.Set("Login", login)
		c.Next()
	}
}
