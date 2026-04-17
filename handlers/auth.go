package handlers

import (
	"karl-s-bar-api/models"
	"karl-s-bar-api/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var users = []models.User{} // 先用内存模拟（后面接 Mongo）

type RegisterRequest struct {
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

    _,err:= h.UserRepository.GetUserByEmail(strings.ToLower(req.Email))
    if (err == nil){
        c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})  
        return      
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user := models.User{
        ID:        time.Now().String(),
        Email:     req.Email,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
    }

    if err := h.UserRepository.InsertUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "user created",
        "email":   user.Email,
    })
}