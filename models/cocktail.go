package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)


type Cocktail struct {
	MongoID      bson.ObjectID `bson:"_id,omitempty" json:"-"`
	ID           string             `bson:"id" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Category     string             `bson:"category" json:"category"`
	Description  string             `bson:"description" json:"description"`
	Image        string             `bson:"image" json:"image"`
	Ingredients  []string           `bson:"ingredients" json:"ingredients"`
	Instructions string             `bson:"instructions" json:"instructions"`
	ABV          int                `bson:"abv" json:"abv"`
	ServingSize  string             `bson:"servingSize" json:"servingSize"`
}