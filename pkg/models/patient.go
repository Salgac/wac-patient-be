package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
	Id             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName      string             `json:"firstName"`
	LastName       string             `json:"lastName"`
	HealthStatuses []HealthStatus     `json:"healthStatuses"`
	Visits         []Visit            `json:"visits"`
}
