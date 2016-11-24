package model

type EventCounter struct {
	Type    string `json:"type" bson:"_id"`
	Counter uint64 `json:"counter" bson:"counter"`
}
