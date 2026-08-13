package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Session struct {
	ID               bson.ObjectID `bson:"_id,omitempty"`
	RefreshTokenHash string        `bson:"refresh_token_hash"`
	Revoked          bool          `bson:"revoked"`
	ExpiredAt        time.Time     `bson:"expires_at"`
	CreatedAt        time.Time     `bson:"created_at"`
}
