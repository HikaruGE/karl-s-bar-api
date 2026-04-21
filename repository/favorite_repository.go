package repository

import (
	"context"
	"karl-s-bar-api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FavoriteRepository interface {
    Create(f *models.Favorite) error
    Exists(userId, cocktailId string) (bool, error)
    GetByUser(userId string) ([]models.Favorite, error)
}

type FavoriteRepositoryMongo struct {
	Collection *mongo.Collection
}

func (r *FavoriteRepositoryMongo) Create(f *models.Favorite) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := r.Collection.InsertOne(ctx, f)
    return err
}

func (r *FavoriteRepositoryMongo) Exists(userId, cocktailId string) (bool, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    count, err := r.Collection.CountDocuments(ctx, bson.M{
        "userId":     userId,
        "cocktailId": cocktailId,
    })

    return count > 0, err
}

func (r *FavoriteRepositoryMongo) GetByUser(userId string) ([]models.Favorite, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := r.Collection.Find(ctx, bson.M{"userId": userId})
    if err != nil {
        return nil, err
    }

    var results []models.Favorite
    if err := cursor.All(ctx, &results); err != nil {
        log.Printf("Failed to decode favorites: %v", err)
        return nil, err
    }

    return results, nil
}