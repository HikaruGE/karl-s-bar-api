package repository

import (
	"context"
	"errors"
	"time"

	"karl-s-bar-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CommentRepository interface {
	CreateComment(comment *models.Comment) error
	GetCommentsByCocktailID(cocktailID bson.ObjectID) ([]models.Comment, error)
	DeleteComment(commentID bson.ObjectID, userID bson.ObjectID) error
}

type CommentRepositoryMongo struct {
	Collection *mongo.Collection
}

func (r *CommentRepositoryMongo) CreateComment(comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(context.Background(), comment)
	return err
}

func (r *CommentRepositoryMongo) GetCommentsByCocktailID(cocktailID bson.ObjectID) ([]models.Comment, error) {
	filter := bson.M{"cocktail_id": cocktailID}
	cursor, err := r.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	if err = cursor.All(context.Background(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepositoryMongo) DeleteComment(commentID bson.ObjectID, userID bson.ObjectID) error {
	filter := bson.M{"_id": commentID, "user_id": userID}
	result, err := r.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("comment not found or not authorized")
	}
	return nil
}