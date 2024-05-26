package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Visit struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Ambulance Ambulance          `json:"ambulance"`
	Timestamp string             `json:"timestamp"`
	Reason    string             `json:"reason"`
	Status    string             `json:"status"`
}
