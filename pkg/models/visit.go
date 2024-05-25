package models

type Visit struct {
	Ambulance Ambulance `json:"ambulance"`
	Timestamp string    `json:"timestamp"`
	Reason    string    `json:"reason"`
}
