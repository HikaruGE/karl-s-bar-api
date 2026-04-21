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
}

type UserRepositoryMongo struct {
	Collection *mongo.Collection
}

func (r *UserRepositoryMongo) GetUserByEmail(email string) (*models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryMongo) InsertUser(user *models.User) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}