package repository

import (
	"karl-s-bar-api/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	MC *mongo.Database
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error){
	collection := r.MC.Collection("users")
	var user *models.User

	err := collection.FindOne(nil, map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

 func (r *UserRepository) InsertUser(user *models.User) error{
	collection := r.MC.Collection("users")
	
	_, err := collection.InsertOne(nil, user)
	if err != nil {
		return err
	}

	return nil
 }