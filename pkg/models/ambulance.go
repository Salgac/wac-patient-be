package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ambulance struct {
	Id   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name"`
}
