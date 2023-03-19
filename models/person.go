package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Qq string
	Mobile int64
}