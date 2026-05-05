package main

import (
	"karl-s-bar-api/handlers"
	middlewares "karl-s-bar-api/middleware"
	"karl-s-bar-api/repository"
	"karl-s-bar-api/validators"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	db, err := repository.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Create database indexes
	if err := repository.CreateIndexes(db); err != nil {
		panic(err)
	}

	cocktailRepo := &repository.CocktailRepositoryMongo{
		Collection: db.Collection("cocktail_recipes"),
	}
	userRepo := &repository.UserRepositoryMongo{
		Collection: db.Collection("users"),
	}
	commentRepo := &repository.CommentRepositoryMongo{
		Collection: db.Collection("comments"),
	}

	healthCheckHandler := &handlers.HealthCheckHandler{}
	cocktailHandler := &handlers.CocktailHandler{
		CocktailRepository: cocktailRepo,
	}
	authHandler := &handlers.AuthHandler{
		UserRepository:    userRepo,
		RegisterValidator: &validators.RegisterValidatorImpl{},
		LoginValidator:    &validators.LoginValidatorImpl{},
	}
	favoriteHandler := &handlers.FavoriteHandler{
		UserRepository:    userRepo,
		FavoriteValidator: &validators.FavoriteValidatorImpl{},
	}
	commentHandler := &handlers.CommentHandler{
		CommentRepository: commentRepo,
		UserRepository:    userRepo,
		CommentValidator:  &validators.CommentValidatorImpl{},
	}

	r.GET("/cheers", healthCheckHandler.HealthCheck)
	r.GET("/cocktails", cocktailHandler.GetCocktailsHandler)
	r.GET("/cocktails/:id", cocktailHandler.GetCocktailByIDHandler)

	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	r.POST("/favorite", middlewares.AuthMiddleware(), favoriteHandler.Create)
	r.GET("/favorite", middlewares.AuthMiddleware(), favoriteHandler.List)
	r.DELETE("/favorite/:cocktailId", middlewares.AuthMiddleware(), favoriteHandler.Delete)

	r.POST("/cocktails/:id/comments", middlewares.AuthMiddleware(), commentHandler.CreateComment)
	r.GET("/cocktails/:id/comments", commentHandler.GetComments)
	r.DELETE("/comments/:commentID", middlewares.AuthMiddleware(), commentHandler.DeleteComment)

	r.Run(":9527")
}
