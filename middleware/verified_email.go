package middlewares

import (
	"net/http"

	"karl-s-bar-api/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func RequireVerifiedEmail(userRepo repository.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        userIDStr := c.GetString("userId")
        if userIDStr == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
            c.Abort()
            return
        }

        userID, err := bson.ObjectIDFromHex(userIDStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
            c.Abort()
            return
        }

        user, err := userRepo.GetUserByID(userID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email status"})
            c.Abort()
            return
        }

        if !user.Verified {
            c.JSON(http.StatusForbidden, gin.H{"error": "email not verified"})
            c.Abort()
            return
        }

        c.Next()
    }
}
