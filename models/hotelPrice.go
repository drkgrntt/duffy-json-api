package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelPrice struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Date    string             `json:"date"`
	Price   float64            `json:"price"`
	Updated time.Time          `json:"updated"`
}
