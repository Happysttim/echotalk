package model

type HashCounter struct {
	ID      string `bson:"_id"`
	Counter int64  `bson:"counter"`
}
