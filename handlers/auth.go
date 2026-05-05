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
    Name     string `json:"name"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthHandler struct {  
	UserRepository       repository.UserRepository
	RegisterValidator    RegisterValidator
	LoginValidator       LoginValidator
}

// RegisterValidator validates register requests
type RegisterValidator interface {
	ValidateRegisterRequest(email, password, name string) error
}

// LoginValidator validates login requests
type LoginValidator interface {
	ValidateLoginRequest(email, password string) error
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Validate register request
	if err := h.RegisterValidator.ValidateRegisterRequest(req.Email, req.Password, req.Name); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user := &models.User{
        Email:     strings.ToLower(req.Email),
        Name:      req.Name,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
    }

    if err := h.UserRepository.InsertUser(user); err != nil {
        // Check if it's a duplicate key error (email already exists)
        if mongo.IsDuplicateKeyError(err) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
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

    // Validate login request
	if err := h.LoginValidator.ValidateLoginRequest(req.Email, req.Password); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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