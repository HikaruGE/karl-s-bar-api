package handlers

import (
	"karl-s-bar-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CocktailHandler struct {
	CocktailGetter CocktailGetter
}

type CocktailGetter interface {
	GetCocktails() ([]models.Cocktail, error)
	GetCocktailByID(id string) (*models.Cocktail, error)
}

func (h *CocktailHandler) GetCocktailsHandler(c *gin.Context) {
	cocktails, err := h.CocktailGetter.GetCocktails()
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
	cocktail, err := h.CocktailGetter.GetCocktailByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cocktail not found",
		})
		return
	}
	c.JSON(http.StatusOK, cocktail)
}
