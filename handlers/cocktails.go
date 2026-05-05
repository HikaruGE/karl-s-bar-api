package handlers

import (
	"karl-s-bar-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CocktailHandler struct {
	CocktailRepository repository.CocktailRepository
}

func (h *CocktailHandler) GetCocktailsHandler(c *gin.Context) {
	cocktails, err := h.CocktailRepository.GetCocktails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch cocktails",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"cocktails": cocktails,
	})
}

func (h *CocktailHandler) GetCocktailByIDHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cocktail ID",
		})
		return
	}
	cocktail, err := h.CocktailRepository.GetCocktailByID(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cocktail not found",
		})
		return
	}
	c.JSON(http.StatusOK, cocktail)
}
