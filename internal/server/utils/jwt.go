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
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.Login, nil
}

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		Authorization := c.Request.Header["Authorization"]
		if len(Authorization) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
			return
		}
		parts := strings.Split(Authorization[0], " ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
			return
		}
		token := parts[1]
		login, err := GetUserLogin(token, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is wrong"})
			return
		}
		c.Set("Login", login)
		c.Next()
	}
}

// Надо доделать
// func RateLimiter(rateLimit int) gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		Authorization := c.Request.Header["Authorization"]
// 		if len(Authorization) != 1 {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
// 			return
// 		}
// 		parts := strings.Split(Authorization[0], " ")
// 		if len(parts) != 2 {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is not set"})
// 			return
// 		}
// 		token := parts[1]
// 		login, err := GetUserLogin(token, secretKey)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is wrong"})
// 			return
// 		}
// 		c.Set("Login", login)
// 		c.Next()
// 	}
// }
