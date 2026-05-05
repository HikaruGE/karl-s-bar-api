package models

import "time"

type FavoriteItem struct {
	CocktailID string    `bson:"cocktailId" json:"cocktailId"`
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
}