package repository

import (
	"karl-s-bar-api/models"

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
	var user *models.User

	err := r.Collection.FindOne(nil, map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

 func (r *UserRepositoryMongo) InsertUser(user *models.User) error{
	_, err := r.Collection.InsertOne(nil, user)
	if err != nil {
		return err
	}

	return nil
 }