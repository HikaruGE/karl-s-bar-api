package models

import "time"

type FavoriteItem struct {
	CocktailID string    `bson:"cocktailId" json:"cocktailId"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}
type Favorite struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"userId" json:"userId"`
	CocktailID string    `bson:"cocktailId" json:"cocktailId"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}
