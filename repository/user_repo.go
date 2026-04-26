package repository

import (
	"context"
	"karl-s-bar-api/models"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(user *models.User) error
	HasFavorite(userID, cocktailID string) (bool, error)
	AddFavorite(userID string, favorite models.FavoriteItem) error
	RemoveFavorite(userID, cocktailID string) error
	GetFavorites(userID string) ([]models.FavoriteItem, error)
}

type UserRepositoryMongo struct {
	Collection *mongo.Collection
}

func (r *UserRepositoryMongo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryMongo) InsertUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepositoryMongo) HasFavorite(userID, cocktailID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	count, err := r.Collection.CountDocuments(ctx, bson.M{
		"_id":                  userObjectID,
		"favorites.cocktailId": cocktailID,
	})
	return count > 0, err
}

func (r *UserRepositoryMongo) AddFavorite(userID string, favorite models.FavoriteItem) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(ctx,
		bson.M{"_id": userObjectID},
		bson.M{"$addToSet": bson.M{"favorites": favorite}},
	)
	return err
}

func (r *UserRepositoryMongo) RemoveFavorite(userID, cocktailID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(ctx,
		bson.M{"_id": userObjectID},
		bson.M{"$pull": bson.M{"favorites": bson.M{"cocktailId": cocktailID}}},
	)
	return err
}

func (r *UserRepositoryMongo) GetFavorites(userID string) ([]models.FavoriteItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.Collection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user.Favorites, nil
}
