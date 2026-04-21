package handlers

import (
	"karl-s-bar-api/models"
	"karl-s-bar-api/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthHandler struct {  
    UserRepository repository.UserRepository
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    user,err:= h.UserRepository.GetUserByEmail(strings.ToLower(req.Email))
    if (err != nil){
        if (err != mongo.ErrNoDocuments) {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})  
            return  
        }
    }
    if user != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})  
        return 
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user = &models.User{
        Email:     strings.ToLower(req.Email),
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
    }

    if err := h.UserRepository.InsertUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "user created",
        "email":   user.Email,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    email := strings.ToLower(req.Email)

    user, err := h.UserRepository.GetUserByEmail(email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
    }

    token, err := generateJWT(user.ID.Hex())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
    })
}

var jwtSecret = []byte("your_secret_key") //todo: move to env var

func generateJWT(userID string) (string, error) {
    claims := jwt.MapClaims{
        "userId": userID,
        "exp":    time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString(jwtSecret)
}