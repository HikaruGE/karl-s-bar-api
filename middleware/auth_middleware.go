package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your_secret_key")

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        authHeader := c.GetHeader("Authorization")

        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        // Bearer xxx
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Get claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
            c.Abort()
            return
        }

        userId, ok := claims["userId"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userId"})
            c.Abort()
            return
        }

        c.Set("userId", userId)

        c.Next()
    }
}