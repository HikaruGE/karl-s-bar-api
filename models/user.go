package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID                     bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	Email                  string         `bson:"email" json:"email"`
	Nickname               string         `bson:"nickname" json:"nickname"`
	Password               string         `bson:"password" json:"-"`
	Verified               bool           `bson:"verified" json:"verified"`
	VerificationToken      string         `bson:"verificationToken,omitempty" json:"-"`
	VerificationTokenExpiry time.Time      `bson:"verificationTokenExpiry,omitempty" json:"-"`
	CreatedAt              time.Time      `bson:"createdAt" json:"createdAt"`
	Favorites              []FavoriteItem `bson:"favorites,omitempty" json:"favorites,omitempty"`
}
