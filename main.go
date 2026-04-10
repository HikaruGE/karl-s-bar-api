package main

import (
	"karl-s-bar-api/data"
	"karl-s-bar-api/handlers"

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

	healthCheckHandler := &handlers.HealthCheckHandler{}
	cocktailHandler := &handlers.CocktailHandler{
		CocktailGetter: &data.CocktailGetterImpl{},
	}

	r.GET("/cheers", healthCheckHandler.HealthCheck)
	r.GET("/cocktails", cocktailHandler.GetCocktailsHandler)
	r.GET("/cocktails/:id", cocktailHandler.GetCocktailByIDHandler)

	r.Run(":9527")
}