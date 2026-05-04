package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CocktailID bson.ObjectID `bson:"cocktail_id" json:"cocktail_id"`
	UserID     bson.ObjectID `bson:"user_id" json:"user_id"`
	UserName   string        `bson:"user_name" json:"user_name"`
	Content    string        `bson:"content" json:"content"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
}