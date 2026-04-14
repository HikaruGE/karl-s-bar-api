package repository

import (
	"context"
	"karl-s-bar-api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CocktailRepo struct {
	MC *mongo.Database
}

func (m *CocktailRepo) GetCocktails() ([]models.Cocktail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.MC.Collection("cocktail_recipes")

	cursor, err := collection.Find(ctx, bson.M{})
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

func (m *CocktailRepo) GetCocktailByID(id string) (*models.Cocktail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.MC.Collection("cocktail_recipes")

	var cocktail *models.Cocktail
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&cocktail)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return cocktail, nil
}