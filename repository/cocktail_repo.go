package repository

import (
	"context"
	"karl-s-bar-api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CocktailRepository interface {
	GetCocktails() ([]models.Cocktail, error)
	GetCocktailByID(id string) (*models.Cocktail, error)
}

type CocktailRepositoryMongo struct {
	Collection *mongo.Collection
}

func (m *CocktailRepositoryMongo) GetCocktails() ([]models.Cocktail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.Collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var cocktails []models.Cocktail
	if err = cursor.All(ctx, &cocktails); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return cocktails, nil
}

func (m *CocktailRepositoryMongo) GetCocktailByID(id string) (*models.Cocktail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var cocktail *models.Cocktail
	err := m.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&cocktail)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return cocktail, nil
}