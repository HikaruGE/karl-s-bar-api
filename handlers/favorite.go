package handlers

import (
	"karl-s-bar-api/models"
	"karl-s-bar-api/repository"
	"time"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
    FavoriteRepository repository.FavoriteRepository
}

type CreateFavoriteRequest struct {
    CocktailID string `json:"cocktailId"`
}

func (h *FavoriteHandler) Create(c *gin.Context) {
    userId := c.GetString("userId")

    var req CreateFavoriteRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }

    exists, _ := h.FavoriteRepository.Exists(userId, req.CocktailID)
    if exists {
        c.JSON(400, gin.H{"error": "already favorited"})
        return
    }

    fav := models.Favorite{
        UserID:     userId,
        CocktailID: req.CocktailID,
        CreatedAt:  time.Now(),
    }

    if err := h.FavoriteRepository.Create(&fav); err != nil {
        c.JSON(500, gin.H{"error": "failed"})
        return
    }

    c.JSON(200, gin.H{"message": "ok"})
}

func (h *FavoriteHandler) List(c *gin.Context) {
    userId := c.GetString("userId")

    favs, err := h.FavoriteRepository.GetByUser(userId)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed"})
        return
    }

    c.JSON(200, favs)
}

func (h *FavoriteHandler) Delete(c *gin.Context) {
    userId := c.GetString("userId")
    cocktailId := c.Param("cocktailId")

    if err := h.FavoriteRepository.DeleteByUserAndCocktail(userId, cocktailId); err != nil {
        c.JSON(500, gin.H{"error": "failed to delete favorite for cocktail " + cocktailId + "and user " + userId})
        return
    }

    c.JSON(200, gin.H{"message": "ok"})
}