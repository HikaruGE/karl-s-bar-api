package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)



type Favorite struct {
    ID         bson.ObjectID    `bson:"_id,omitempty" json:"id"`
    UserID     string    `bson:"userId" json:"userId"`
    CocktailID string    `bson:"cocktailId" json:"cocktailId"`
    CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}