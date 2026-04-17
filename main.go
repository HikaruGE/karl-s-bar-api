package main

import (
	"karl-s-bar-api/handlers"
	"karl-s-bar-api/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(func (c *gin.Context)  {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return	
		}
		
		c.Next()
	})

	db := repository.ConnectDB()

	cocktailRepo := &repository.CocktailRepositoryMongo{
		Collection: db.Collection("cocktail_recipes"),
	}
	userRepo := &repository.UserRepositoryMongo{
		Collection: db.Collection("users"),
	}

	healthCheckHandler := &handlers.HealthCheckHandler{}
	cocktailHandler := &handlers.CocktailHandler{
		CocktailRepository: cocktailRepo,
	}
	authHandler := &handlers.AuthHandler{UserRepository: userRepo}

	r.GET("/cheers", healthCheckHandler.HealthCheck)
	r.GET("/cocktails", cocktailHandler.GetCocktailsHandler)
	r.GET("/cocktails/:id", cocktailHandler.GetCocktailByIDHandler)

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	r.Run(":9527")
}