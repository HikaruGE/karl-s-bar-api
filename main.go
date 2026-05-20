package main

import (
	"karl-s-bar-api/handlers"
	"karl-s-bar-api/mail"
	middlewares "karl-s-bar-api/middleware"
	"karl-s-bar-api/repository"
	"karl-s-bar-api/utils"
	"karl-s-bar-api/validators"
	"os"
	"strconv"
	"strings"

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

	smtpPort := 587
	if portStr := strings.TrimSpace(os.Getenv("SMTP_PORT")); portStr != "" {
		if parsed, err := strconv.Atoi(portStr); err == nil {
			smtpPort = parsed
		}
	}

	smtpHost := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	smtpUser := strings.TrimSpace(os.Getenv("SMTP_USER"))
	smtpPass := strings.ReplaceAll(strings.TrimSpace(os.Getenv("SMTP_PASS")), " ", "")
	mailFrom := strings.TrimSpace(os.Getenv("MAIL_FROM"))

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || mailFrom == "" {
		panic("SMTP_HOST, SMTP_USER, SMTP_PASS, and MAIL_FROM must be set")
	}

	emailSender := mail.NewSMTPEmailSender(
		smtpHost,
		smtpPort,
		smtpUser,
		smtpPass,
		mailFrom,
	)

	appBaseURL := os.Getenv("APP_BASE_URL")
	if appBaseURL == "" {
		appBaseURL = "http://localhost:9527"
	}

	tokenGenerator := utils.NewRandomTokenGenerator(32)

	healthCheckHandler := &handlers.HealthCheckHandler{}
	cocktailHandler := &handlers.CocktailHandler{
		CocktailRepository: cocktailRepo,
	}
	authHandler := &handlers.AuthHandler{
		UserRepository:    userRepo,
		RegisterValidator: &validators.RegisterValidatorImpl{},
		LoginValidator:    &validators.LoginValidatorImpl{},
		EmailSender:       emailSender,
		TokenGenerator:    tokenGenerator,
		AppBaseURL:       appBaseURL,
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
	r.GET("/auth/verify", authHandler.VerifyEmail)
	r.POST("/auth/resend-verification", authHandler.ResendVerification)
	r.GET("/auth/profile", middlewares.AuthMiddleware(), middlewares.RequireVerifiedEmail(userRepo), authHandler.Profile)

	r.POST("/favorite", middlewares.AuthMiddleware(), middlewares.RequireVerifiedEmail(userRepo), favoriteHandler.Create)
	r.GET("/favorite", middlewares.AuthMiddleware(), middlewares.RequireVerifiedEmail(userRepo), favoriteHandler.List)
	r.DELETE("/favorite/:cocktailId", middlewares.AuthMiddleware(), middlewares.RequireVerifiedEmail(userRepo), favoriteHandler.Delete)

	r.POST("/cocktails/:id/comments", middlewares.AuthMiddleware(), middlewares.RequireVerifiedEmail(userRepo), commentHandler.CreateComment)
	r.GET("/cocktails/:id/comments", commentHandler.GetComments)
	r.DELETE("/comments/:commentID", middlewares.AuthMiddleware(), middlewares.RequireVerifiedEmail(userRepo), commentHandler.DeleteComment)

	r.Run(":9527")
}
