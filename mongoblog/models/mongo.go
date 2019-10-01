package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mongo stands for mongo BSON ObjectID type.
type Mongo struct {
	OID primitive.ObjectID `bson:"_id,omitempty" json:"-"`
}
