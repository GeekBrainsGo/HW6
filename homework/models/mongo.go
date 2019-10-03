package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mongo struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}

func (m *Mongo) GetMongoCollectionName() string {
	panic("GetMongoCollectionName not implemented")
	return ""
}
