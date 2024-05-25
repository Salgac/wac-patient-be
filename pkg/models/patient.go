package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
	Id               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName        string             `json:"firstName"`
	LastName         string             `json:"lastName"`
	Age              int                `json:"age"`
	HealthConditions []HealthCondition  `json:"healthConditions"`
	Visits           []Visit            `json:"visits"`
}
