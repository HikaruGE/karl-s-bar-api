package handlers

import (
	"fmt"
	"karl-s-bar-api/mail"
	"karl-s-bar-api/models"
	"karl-s-bar-api/repository"
	"karl-s-bar-api/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Nickname string `json:"nickname"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type ResendVerificationRequest struct {
    Email string `json:"email"`
}

type AuthHandler struct {
	UserRepository  repository.UserRepository
	RegisterValidator RegisterValidator
	LoginValidator  LoginValidator
	EmailSender     mail.EmailSender
	TokenGenerator  utils.TokenGenerator
	AppBaseURL      string
}

// RegisterValidator validates register requests
type RegisterValidator interface {
    ValidateRegisterRequest(email, password, nickname string) error
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

    if err := h.RegisterValidator.ValidateRegisterRequest(req.Email, req.Password, req.Nickname); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    verificationToken, err := h.TokenGenerator.GenerateToken()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification token"})
        return
    }

    user := &models.User{
        Email:                  strings.ToLower(req.Email),
        Nickname:               req.Nickname,
        Password:               string(hashedPassword),
        Verified:               false,
        VerificationToken:      verificationToken,
        VerificationTokenExpiry: time.Now().Add(24 * time.Hour),
        CreatedAt:              time.Now(),
    }

    if err := h.UserRepository.InsertUser(user); err != nil {
        if mongo.IsDuplicateKeyError(err) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
    }

    verificationURL := fmt.Sprintf("%s/auth/verify?token=%s", strings.TrimRight(h.AppBaseURL, "/"), verificationToken)
    if err := h.EmailSender.SendVerificationEmail(user.Email, user.Nickname, verificationURL); err != nil {
        fmt.Printf("send verification email failed: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification email", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "user created, verification email sent",
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

    if !user.Verified {
        c.JSON(http.StatusForbidden, gin.H{"error": "email not verified"})
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

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
    token := c.Query("token")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "missing token"})
        return
    }

    user, err := h.UserRepository.GetUserByVerificationToken(token)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify token"})
        return
    }

    if user.Verified {
        c.JSON(http.StatusOK, gin.H{"message": "email already verified"})
        return
    }

    if user.VerificationTokenExpiry.Before(time.Now()) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "verification token expired"})
        return
    }

    if err := h.UserRepository.MarkUserVerified(user.ID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
    var req ResendVerificationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    email := strings.ToLower(req.Email)
    user, err := h.UserRepository.GetUserByEmail(email)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
        return
    }

    if user.Verified {
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already verified"})
        return
    }

    verificationToken, err := h.TokenGenerator.GenerateToken()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification token"})
        return
    }

    expiry := time.Now().Add(24 * time.Hour)
    if err := h.UserRepository.UpdateVerificationToken(user.ID, verificationToken, expiry); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update verification token"})
        return
    }

    verificationURL := fmt.Sprintf("%s/auth/verify?token=%s", strings.TrimRight(h.AppBaseURL, "/"), verificationToken)
    if err := h.EmailSender.SendVerificationEmail(user.Email, user.Nickname, verificationURL); err != nil {
        fmt.Printf("resend verification email failed: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification email", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "verification email resent"})
}

func (h *AuthHandler) Profile(c *gin.Context) {
    userIDStr := c.GetString("userId")
    if userIDStr == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    userID, err := bson.ObjectIDFromHex(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

    user, err := h.UserRepository.GetUserByID(userID)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "email":    user.Email,
        "nickname": user.Nickname,
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