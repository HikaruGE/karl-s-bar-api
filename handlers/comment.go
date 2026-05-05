package handlers

import (
	"net/http"
	"time"

	"karl-s-bar-api/models"
	"karl-s-bar-api/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CommentHandler struct {
	CommentRepository  repository.CommentRepository
	UserRepository     repository.UserRepository
	CommentValidator   CommentValidator
}

type CommentValidator interface {
	ValidateCreateCommentRequest(content string) error
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	cocktailIDStr := c.Param("id")
	cocktailID, err := bson.ObjectIDFromHex(cocktailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cocktail ID"})
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate comment request
	if err := h.CommentValidator.ValidateCreateCommentRequest(req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from JWT
	userIDStr := c.GetString("userId")
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user details from database
	user, err := h.UserRepository.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	comment := &models.Comment{
		CocktailID: cocktailID,
		UserID:     userID,
		UserName:   user.Email, // Use email as username
		Content:    req.Content,
		CreatedAt:  time.Now(),
	}

	if err := h.CommentRepository.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	cocktailIDStr := c.Param("id")
	cocktailID, err := bson.ObjectIDFromHex(cocktailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cocktail ID"})
		return
	}

	comments, err := h.CommentRepository.GetCommentsByCocktailID(cocktailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIDStr := c.Param("commentID")
	commentID, err := bson.ObjectIDFromHex(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Get user info from JWT
	userIDStr := c.GetString("userId")
	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.CommentRepository.DeleteComment(commentID, userID); err != nil {
		if err.Error() == "comment not found or not authorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}