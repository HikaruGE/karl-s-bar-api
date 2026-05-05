package handlers

import (
	"karl-s-bar-api/models"
	"karl-s-bar-api/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	UserRepository    repository.UserRepository
	FavoriteValidator FavoriteValidator
}

type FavoriteValidator interface {
	ValidateCreateFavoriteRequest(cocktailID string) error
	ValidateDeleteFavoriteRequest(cocktailID string) error
}

type CreateFavoriteRequest struct {
	CocktailID string `json:"cocktailId"`
}

func (h *FavoriteHandler) Create(c *gin.Context) {
	userId := c.GetString("userId")

	var req CreateFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate favorite request
	if err := h.FavoriteValidator.ValidateCreateFavoriteRequest(req.CocktailID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := h.UserRepository.HasFavorite(userId, req.CocktailID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already favorited"})
		return
	}

	fav := models.FavoriteItem{
		CocktailID: req.CocktailID,
		CreatedAt:  time.Now(),
	}

	if err := h.UserRepository.AddFavorite(userId, fav); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *FavoriteHandler) List(c *gin.Context) {
	userId := c.GetString("userId")

	favs, err := h.UserRepository.GetFavorites(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, favs)
}

func (h *FavoriteHandler) Delete(c *gin.Context) {
	userId := c.GetString("userId")
	cocktailId := c.Param("cocktailId")

	// Validate cocktailId
	if err := h.FavoriteValidator.ValidateDeleteFavoriteRequest(cocktailId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UserRepository.RemoveFavorite(userId, cocktailId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
